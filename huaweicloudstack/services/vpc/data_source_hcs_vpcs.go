package vpc

import (
	"context"
	v3Vpcs "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v3/vpcs"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/vpcs"
)

func DataSourceVpcs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpcs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secondary_cidrs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceVpcsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	region := hcsConfig.GetRegion(d)
	client, err := hcsConfig.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	listOpts := vpcs.ListOpts{
		ID:                  d.Get("id").(string),
		Name:                d.Get("name").(string),
		Status:              d.Get("status").(string),
		CIDR:                d.Get("cidr").(string),
		EnterpriseProjectID: hcsConfig.DataGetEnterpriseProjectID(d),
	}

	vpcList, err := vpcs.List(client, listOpts)
	if err != nil {
		return diag.Errorf("unable to retrieve vpcs: %s", err)
	}

	log.Printf("[DEBUG] Retrieved Vpc using given filter: %+v", vpcList)

	var vpcs []map[string]interface{}
	var ids []string
	for _, vpcResource := range vpcList {
		vpc := map[string]interface{}{
			"id":                    vpcResource.ID,
			"name":                  vpcResource.Name,
			"cidr":                  vpcResource.CIDR,
			"enterprise_project_id": vpcResource.EnterpriseProjectID,
			"status":                vpcResource.Status,
			"description":           vpcResource.Description,
		}

		vpcs = append(vpcs, vpc)
		ids = append(ids, vpcResource.ID)

		// save VirtualPrivateCloudV3 extend_cidr
		vpcV3Client, v3Err := hcsConfig.NetworkingV3Client(hcsConfig.GetRegion(d))
		if v3Err != nil {
			return diag.Errorf("error creating VPC v3 client: %s", err)
		}

		res, err := v3Vpcs.Get(vpcV3Client, vpcResource.ID).Extract()
		if err != nil {
			diag.Errorf("error retrieving VPC (%s) v3 detail: %s", d.Id(), err)
		}
		vpc["secondary_cidrs"] = res.ExtendCidrs

		vpcs = append(vpcs, vpc)
		ids = append(ids, vpcResource.ID)
	}
	log.Printf("[DEBUG] Vpc List after filter, count=%d :%+v", len(vpcs), vpcs)

	mErr := d.Set("vpcs", vpcs)
	if mErr != nil {
		return diag.Errorf("set vpcs err:%s", mErr)
	}

	d.SetId(hashcode.Strings(ids))
	return nil
}
