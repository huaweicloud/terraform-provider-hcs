package vpc

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v2/dc/virtual_gateways"
	"log"
)

func DataSourceVirtualGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVirtualGatewayRead,

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
			"vpc_group": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_ep_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_ep_group": {
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

func dataSourceVirtualGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	vgwClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Virtual Gateway client: %s", err)
	}

	listOpts := virtual_gateways.ListOpts{
		ID:   d.Get("id").(string),
		Name: d.Get("name").(string),
	}

	n, err := virtual_gateways.List(vgwClient, listOpts)
	if len(n) <= 0 || err != nil {
		return diag.Errorf("error querying Virtual Gateway: %s", err)
	}
	log.Printf("[INFO] Retrieved Virtual Gateway using given filter %s: %+v", n[0].ID, n[0])

	vpcGroups, err := getVpcGroups(n[0].VpcGroup, vgwClient)
	if err != nil {
		return diag.Errorf("error getting Virtual Gateway vpc groups: %s", err)
	}

	d.SetId(n[0].ID)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", n[0].Name),
		d.Set("status", n[0].Status),
		d.Set("description", n[0].Description),
		d.Set("vpc_group", vpcGroups),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Virtual Gateway fields: %s", err)
	}

	return nil
}
