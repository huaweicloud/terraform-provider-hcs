package vpc

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v2/dc/dc_endpoint_groups"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v2/dc/virtual_interfaces"
	"log"
)

func DataSourceVirtualInterface() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVirtualInterfaceRead,

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
			"direct_connect_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vgw_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remote_ep_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remote_ep_group": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"link_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interface_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hosting_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_gateway_v4_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_gateway_v6_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_gateway_v4_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_gateway_v6_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vlan": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bgp_asn": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bgp_asn_dot": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceVirtualInterfaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	vifClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Virtual Interface client: %s", err)
	}

	listOpts := virtual_interfaces.ListOpts{
		ID:   d.Get("id").(string),
		Name: d.Get("name").(string),
	}

	n, err := virtual_interfaces.List(vifClient, listOpts)
	if len(n) <= 0 || err != nil {
		return diag.Errorf("error querying Virtual Interface: %s", err)
	}
	log.Printf("[INFO] Retrieved Virtual Interaface using given filter %s: %+v", n[0].ID, n[0])

	endpointGroup, err := dc_endpoint_groups.Get(vifClient, n[0].RemoteEpGroupId).Extract()
	if err != nil {
		return diag.Errorf("error retrieving Dc Endpoint Group of Virtual Interface: %s", err)
	}

	d.SetId(n[0].ID)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", n[0].Name),
		d.Set("status", n[0].Status),
		d.Set("description", n[0].Description),
		d.Set("direct_connect_id", n[0].DirectConnectId),
		d.Set("vgw_id", n[0].VgwId),
		d.Set("remote_ep_group_id", n[0].RemoteEpGroupId),
		d.Set("remote_ep_group", endpointGroup.Endpoints),
		d.Set("link_infos", getLinkInfos(&n[0])),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Virtual Interface fields: %s", err)
	}

	return nil
}
