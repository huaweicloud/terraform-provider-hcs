package vpc

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v2/dc/dc_endpoint_groups"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v2/dc/virtual_gateways"
	"log"
	"sort"
)

func ResourceVirtualGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVirtualGatewayCreate,
		ReadContext:   resourceVirtualGatewayRead,
		UpdateContext: resourceVirtualGatewayUpdate,
		DeleteContext: resourceVritualGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_group": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"local_ep_group_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"local_ep_group": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.IsCIDR,
							},
							Set: schema.HashString,
						},
					},
				},
				Set: hashVpcGroup,
			},
		},
	}
}

func resourceVirtualGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	vgwClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Virtual Gateway client: %s", err)
	}

	name := d.Get("name").(string)
	vpcGroupData := d.Get("vpc_group").(*schema.Set).List()
	vpcGroups, err := createVgwVpcGroups(name, vpcGroupData, vgwClient)
	if err != nil {
		return diag.Errorf("error creating Virtual Gateway vpc group: %s", err)
	}

	defer func() {
		// delete vpc groups when encounter errors
		if err != nil && len(vpcGroups) > 0 {
			for _, item := range vpcGroups {
				deleteErr := dc_endpoint_groups.Delete(vgwClient, item.LocalEpGroupId).Err
				if deleteErr != nil {
					log.Printf("[ERROR] error deleting Dc Endpoint Group of VGW: %#v, error: %s", item, deleteErr)
				}
			}
		}
	}()

	createOpts := virtual_gateways.CreateOpts{
		Name:        name,
		Description: d.Get("description").(string),
		VpcGroup:    vpcGroups,
	}

	n, err := virtual_gateways.Create(vgwClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating Virtual Gateway: %s", err)
	}

	d.SetId(n.ID)
	log.Printf("[DEBUG] Virtual Gateway ID: %s", n.ID)
	return resourceVirtualGatewayRead(ctx, d, meta)
}

func resourceVirtualGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	vgwClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Virtual Gateway client: %s", err)
	}

	n, err := virtual_gateways.Get(vgwClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Virtual Gateway Connection")
	}

	vpcGroups, err := getVpcGroups(n.VpcGroup, vgwClient)
	if err != nil {
		return diag.Errorf("error getting Virtual Gateway vpc groups: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", n.Name),
		d.Set("status", n.Status),
		d.Set("description", n.Description),
		d.Set("vpc_group", vpcGroups),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Virtual Gateway fields: %s", err)
	}

	return nil
}

func resourceVirtualGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	vgwClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Virtual Gateway client: %s", err)
	}

	name := d.Get("name").(string)

	updateOpts := virtual_gateways.UpdateOpts{
		Name:        name,
		Description: d.Get("description").(string),
	}

	if d.HasChange("vpc_group") {
		var vpcGroupOptsList []virtual_gateways.VpcGroupOpts
		o, n := d.GetChange("vpc_group")
		toDelVpcGroup := o.(*schema.Set).Difference(n.(*schema.Set)).List()
		toAddVpcGroup := n.(*schema.Set).Difference(o.(*schema.Set)).List()
		unChangedVpcGroup := n.(*schema.Set).Intersection(o.(*schema.Set)).List()

		if len(toAddVpcGroup) > 0 {
			vpcGroups, err := createVgwVpcGroups(name, toAddVpcGroup, vgwClient)
			if err != nil {
				return diag.Errorf("error creating Virtual Gateway vpc group: %s", err)
			}
			defer func() {
				// delete new vpc groups when encounter errors
				if err != nil && len(vpcGroups) > 0 {
					for _, item := range vpcGroups {
						deleteErr := dc_endpoint_groups.Delete(vgwClient, item.LocalEpGroupId).Err
						if deleteErr != nil {
							log.Printf("[ERROR] error deleting Dc Endpoint Group of VGW: %#v, error: %s", item, deleteErr)
						}
					}
				}
			}()
			vpcGroupOptsList = append(vpcGroupOptsList, vpcGroups...)
		}

		for _, item := range unChangedVpcGroup {
			vpcGroup := item.(map[string]interface{})
			vpcGroupOpts := virtual_gateways.VpcGroupOpts{
				VpcId:          vpcGroup["vpc_id"].(string),
				LocalEpGroupId: vpcGroup["local_ep_group_id"].(string),
			}
			vpcGroupOptsList = append(vpcGroupOptsList, vpcGroupOpts)
		}
		updateOpts.VpcGroup = vpcGroupOptsList

		defer func() {
			// delete old vpc groups
			if err == nil && len(toDelVpcGroup) > 0 {
				for _, item := range toDelVpcGroup {
					vpcGroup := item.(map[string]interface{})
					localEpGroupId := vpcGroup["local_ep_group_id"].(string)
					deleteErr := dc_endpoint_groups.Delete(vgwClient, localEpGroupId).Err
					if deleteErr != nil {
						log.Printf("[ERROR] error deleting Dc Endpoint Group of VGW: %#v, error: %s", item, deleteErr)
					}
				}
			}
		}()
	}

	_, err = virtual_gateways.Update(vgwClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating Virtual Gateway: %s", err)
	}

	return resourceVirtualGatewayRead(ctx, d, meta)
}

func resourceVritualGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	vgwClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Virtual Gateway client: %s", err)
	}

	err = virtual_gateways.Delete(vgwClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting Virtual Gateway: %s", err)
	}

	vpcGroupData := d.Get("vpc_group").(*schema.Set).List()
	if len(vpcGroupData) > 0 {
		for _, item := range vpcGroupData {
			vpcGroup := item.(map[string]interface{})
			localEpGroupId := vpcGroup["local_ep_group_id"].(string)
			deleteErr := dc_endpoint_groups.Delete(vgwClient, localEpGroupId).Err
			if deleteErr != nil {
				log.Printf("[ERROR] error deleting VGW endpint group: %#v, error: %s", vpcGroup, deleteErr)
			}
		}
	}

	return nil
}

func createVgwVpcGroups(name string, vpcGroups []interface{}, client *golangsdk.ServiceClient) ([]virtual_gateways.VpcGroupOpts, error) {
	var vpcGroupOptsList []virtual_gateways.VpcGroupOpts

	for _, item := range vpcGroups {
		vpcGroup := item.(map[string]interface{})
		vpcId := vpcGroup["vpc_id"].(string)
		epGroupList := vpcGroup["local_ep_group"].(*schema.Set).List()
		localEpGroup := make([]string, len(epGroupList))
		for i, v := range epGroupList {
			localEpGroup[i] = v.(string)
		}

		createOpts := dc_endpoint_groups.CreateOpts{
			Name:      name,
			Type:      "cidr",
			Endpoints: localEpGroup,
		}

		log.Printf("[DEBUG] Create Dc Endpoint Group of VGW: %#v", createOpts)
		dcEndpointGroup, err := dc_endpoint_groups.Create(client, createOpts).Extract()
		if err != nil {
			return vpcGroupOptsList, err
		}
		log.Printf("[DEBUG] Dc Endpoint Group of VGW created: %#v", dcEndpointGroup)

		vpcGroupOpts := virtual_gateways.VpcGroupOpts{
			VpcId:          vpcId,
			LocalEpGroupId: dcEndpointGroup.ID,
		}
		vpcGroupOptsList = append(vpcGroupOptsList, vpcGroupOpts)
	}

	return vpcGroupOptsList, nil
}

func getVpcGroups(vpcGroupResults []virtual_gateways.VpcGroup, client *golangsdk.ServiceClient) ([]map[string]interface{}, error) {
	var vpcGroups []map[string]interface{}
	for _, vpcGroup := range vpcGroupResults {
		dcEndpointGroup, err := dc_endpoint_groups.Get(client, vpcGroup.LocalEpGroupId).Extract()
		if err != nil {
			return vpcGroups, err
		}
		endpoints := dcEndpointGroup.Endpoints
		vpcGroups = append(vpcGroups, map[string]interface{}{
			"vpc_id":            vpcGroup.VpcId,
			"local_ep_group_id": vpcGroup.LocalEpGroupId,
			"local_ep_group":    endpoints,
		})
	}
	return vpcGroups, nil
}

func hashVpcGroup(v interface{}) int {
	vpcGroup := v.(map[string]interface{})
	vpcId := vpcGroup["vpc_id"].(string)
	epGroupList := vpcGroup["local_ep_group"].(*schema.Set).List()
	localEpGroup := make([]string, len(epGroupList))
	for i, epGroup := range epGroupList {
		localEpGroup[i] = epGroup.(string)
	}
	sort.Strings(localEpGroup)
	return schema.HashString(fmt.Sprintf("%s-%s", vpcId, localEpGroup))
}
