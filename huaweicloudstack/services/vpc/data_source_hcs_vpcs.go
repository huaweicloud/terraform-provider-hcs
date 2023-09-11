package vpc

import (
	"context"
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
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
			"id":          vpcResource.ID,
			"name":        vpcResource.Name,
			"cidr":        vpcResource.CIDR,
			"status":      vpcResource.Status,
			"description": vpcResource.Description,
		}

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
