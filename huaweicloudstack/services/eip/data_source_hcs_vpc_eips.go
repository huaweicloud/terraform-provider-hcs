package eip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/eips"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/logp"
)

func DataSourceVpcEips() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcEipsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"eips": {
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bandwidth_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_share_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVpcEipsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := config.GetHcsConfig(meta)
	region := config.GetRegion(d)
	client, err := config.NetworkingV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud Networking client: %s", err)
	}

	listOpts := &eips.ListOpts{
		EnterpriseProjectId: config.DataGetEnterpriseProjectID(d),
	}

	pages, err := eips.List(client, listOpts).AllPages()
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve eips: %s ", err)
	}

	allEips, err := eips.ExtractPublicIPs(pages)
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve eips: %s ", err)
	}

	logp.Printf("[DEBUG] Retrieved eips using given filter: %+v", allEips)

	var eips []map[string]interface{}
	var ids []string
	for _, item := range allEips {
		eip := map[string]interface{}{
			"id":                    item.ID,
			"status":                NormalizeEipStatus(item.Status),
			"type":                  item.Type,
			"private_ip":            item.PrivateAddress,
			"public_ip":             item.PublicAddress,
			"port_id":               item.PortID,
			"enterprise_project_id": item.EnterpriseProjectID,
			"bandwidth_id":          item.BandwidthID,
			"bandwidth_size":        item.BandwidthSize,
			"bandwidth_name":        item.BandwidthName,
			"bandwidth_share_type":  item.BandwidthShareType,
		}

		eips = append(eips, eip)
		ids = append(ids, item.ID)
	}
	logp.Printf("[DEBUG]Eips List after filter, count=%d :%+v", len(eips), eips)

	mErr := d.Set("eips", eips)
	if mErr != nil {
		return fmtp.DiagErrorf("set eips err:%s", mErr)
	}

	d.SetId(hashcode.Strings(ids))
	return nil
}
