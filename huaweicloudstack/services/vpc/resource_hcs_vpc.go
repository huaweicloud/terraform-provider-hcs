package vpc

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/vpcs"
	v3Vpcs "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v3/vpcs"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func ResourceVirtualPrivateCloudV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVirtualPrivateCloudCreate,
		ReadContext:   resourceVirtualPrivateCloudRead,
		UpdateContext: resourceVirtualPrivateCloudUpdate,
		DeleteContext: resourceVirtualPrivateCloudDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ // request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 64),
					validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa50-9a-zA-Z-_\\.]*$"),
						"only letters, digits, underscores (_), hyphens (-), and dot (.) are allowed"),
				),
			},
			"cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.ValidateCIDR,
			},
			"secondary_cidrs": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: utils.ValidateCIDR,
				},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"routes": {
				Type:       schema.TypeList,
				Computed:   true,
				Deprecated: "use hcs_vpc_route_table data source to get all routes",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nexthop": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceVirtualPrivateCloudCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := config.GetHcsConfig(meta)
	region := config.GetRegion(d)
	vpcClient, err := config.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	createOpts := vpcs.CreateOpts{
		Name: d.Get("name").(string),
		CIDR: d.Get("cidr").(string),
	}

	epsID := common.GetEnterpriseProjectID(d, config)
	if epsID != "" {
		createOpts.EnterpriseProjectID = epsID
	}

	n, err := vpcs.Create(vpcClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating VPC: %s", err)
	}

	d.SetId(n.ID)
	log.Printf("[DEBUG] Vpc ID: %s", n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForVpcActive(vpcClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return diag.Errorf(
			"error waiting for Vpc (%s) to become ACTIVE: %s",
			n.ID, stateErr)
	}
	var extendCidrs []string
	if v, ok := d.GetOk("secondary_cidrs"); ok {
		extendCidrs = utils.ExpandToStringList(v.(*schema.Set).List())
	}

	if len(extendCidrs) > 0 {
		vpcV3Client, v3Err := config.NetworkingV3Client(region)
		if v3Err != nil {
			return diag.Errorf("error creating VPC v3 client: %s", err)
		}
		updatOpts := v3Vpcs.UpdateOpts{
			ExtendCidrs: extendCidrs,
		}
		_, newErr := v3Vpcs.AddSecondaryCIDR(vpcV3Client, d.Id(), updatOpts).Extract()

		if newErr != nil {
			return diag.Errorf("error adding VPC secondary CIDRs: %s", newErr)
		}
	}

	return resourceVirtualPrivateCloudRead(ctx, d, meta)
}

// GetVpcById is a method to obtain vpc informations from special region through vpc ID.
func GetVpcById(config *config.HcsConfig, region, vpcId string) (*vpcs.Vpc, error) {
	client, err := config.NetworkingV1Client(region)
	if err != nil {
		return nil, err
	}

	return vpcs.Get(client, vpcId).Extract()
}

func resourceVirtualPrivateCloudRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := config.GetHcsConfig(meta)
	n, err := GetVpcById(config, config.GetRegion(d), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error obtain VPC information")
	}

	d.Set("name", n.Name)
	d.Set("cidr", n.CIDR)
	d.Set("status", n.Status)
	d.Set("enterprise_project_id", n.EnterpriseProjectID)
	d.Set("region", config.GetRegion(d))

	// save route tables
	routes := make([]map[string]interface{}, len(n.Routes))
	for i, rtb := range n.Routes {
		route := map[string]interface{}{
			"destination": rtb.DestinationCIDR,
			"nexthop":     rtb.NextHop,
		}
		routes[i] = route
	}
	d.Set("routes", routes)
	vpcV3Client, v3Err := config.NetworkingV3Client(config.GetRegion(d))
	if v3Err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}
	res, err := v3Vpcs.Get(vpcV3Client, d.Id()).Extract()
	if err != nil {
		diag.Errorf("error retrieving VPC (%s) v3 detail: %s", d.Id(), err)
	}
	d.Set("secondary_cidrs", res.ExtendCidrs)

	return nil
}

func resourceVirtualPrivateCloudUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := config.GetHcsConfig(meta)
	region := config.GetRegion(d)
	vpcClient, err := config.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	vpcID := d.Id()
	if d.HasChanges("name", "cidr") {
		updateOpts := vpcs.UpdateOpts{
			Name: d.Get("name").(string),
			CIDR: d.Get("cidr").(string),
		}

		_, err = vpcs.Update(vpcClient, vpcID, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating VPC: %s", err)
		}
	}
	vpcV3Client, v3Err := config.NetworkingV3Client(region)
	if v3Err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	if d.HasChanges("secondary_cidrs") {
		oldRaws, newRaws := d.GetChange("secondary_cidrs")
		needRemovecidrs := utils.ExpandToStringListBySet(oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set)))
		newExtendCidrs := utils.ExpandToStringListBySet(newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set)))
		if len(needRemovecidrs) > 0 {
			removeOpts := v3Vpcs.UpdateOpts{
				ExtendCidrs: needRemovecidrs,
			}
			_, removeErr := v3Vpcs.RemoveSecondaryCIDR(vpcV3Client, d.Id(), removeOpts).Extract()

			if removeErr != nil {
				return diag.Errorf("error deleting VPC secondary CIDRs: %s", err)
			}
		}
		if len(newExtendCidrs) > 0 {
			addOpts := v3Vpcs.UpdateOpts{
				ExtendCidrs: newExtendCidrs,
			}
			_, addErr := v3Vpcs.AddSecondaryCIDR(vpcV3Client, d.Id(), addOpts).Extract()

			if addErr != nil {
				return diag.Errorf("error adding VPC secondary CIDRs: %s", addErr)
			}
		}
	}

	return resourceVirtualPrivateCloudRead(ctx, d, meta)
}

func resourceVirtualPrivateCloudDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	config := config.GetHcsConfig(meta)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForVpcDelete(vpcClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting VPC %s: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func waitForVpcActive(vpcClient *golangsdk.ServiceClient, vpcId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := vpcs.Get(vpcClient, vpcId).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.Status == "OK" {
			return n, "ACTIVE", nil
		}

		// If vpc status is other than Ok, send error
		if n.Status == "DOWN" {
			return nil, "", fmt.Errorf("VPC status: '%s'", n.Status)
		}

		return n, n.Status, nil
	}
}

func waitForVpcDelete(vpcClient *golangsdk.ServiceClient, vpcId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := vpcs.Get(vpcClient, vpcId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully delete VPC %s", vpcId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = vpcs.Delete(vpcClient, vpcId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully delete VPC %s", vpcId)
				return r, "DELETED", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "ACTIVE", nil
				}
			}
			return r, "ACTIVE", err
		}

		return r, "ACTIVE", nil
	}
}
