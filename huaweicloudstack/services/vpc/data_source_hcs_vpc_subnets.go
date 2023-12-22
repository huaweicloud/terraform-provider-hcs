package vpc

import (
	"context"
	"log"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/subnets"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceVpcSubnets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcSubnetsRead,

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
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateway_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"primary_dns": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secondary_dns": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnets": {
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
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"primary_dns": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secondary_dns": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dhcp_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"dns_list": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "schema: Deprecated",
						},
						"ipv4_subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6_subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ipv6_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6_gateway": {
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

func dataSourceVpcSubnetsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	region := hcsConfig.GetRegion(d)
	client, err := hcsConfig.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	if err != nil {
		return diag.Errorf("error creating VPC V2 client: %s", err)
	}

	listOpts := subnets.ListOpts{
		ID:            d.Get("id").(string),
		Name:          d.Get("name").(string),
		CIDR:          d.Get("cidr").(string),
		Status:        d.Get("status").(string),
		GatewayIP:     d.Get("gateway_ip").(string),
		PRIMARY_DNS:   d.Get("primary_dns").(string),
		SECONDARY_DNS: d.Get("secondary_dns").(string),
		VPC_ID:        d.Get("vpc_id").(string),
	}

	subnetList, err := subnets.List(client, listOpts)
	if err != nil {
		return diag.Errorf("unable to retrieve subnets: %s", err)
	}

	log.Printf("[DEBUG] Retrieved subnets using given filter: %+v", subnetList)

	var subnets []map[string]interface{}
	var ids []string
	for _, item := range subnetList {
		subnet := map[string]interface{}{
			"id":             item.ID,
			"name":           item.Name,
			"description":    item.Description,
			"cidr":           item.CIDR,
			"status":         item.Status,
			"gateway_ip":     item.GatewayIP,
			"dns_list":       item.DnsList,
			"ipv6_enable":    item.EnableIPv6,
			"dhcp_enable":    item.EnableDHCP,
			"primary_dns":    item.PRIMARY_DNS,
			"secondary_dns":  item.SECONDARY_DNS,
			"vpc_id":         item.VPC_ID,
			"subnet_id":      item.SubnetId,
			"ipv4_subnet_id": item.SubnetId,
			"ipv6_subnet_id": item.IPv6SubnetId,
			"ipv6_cidr":      item.IPv6CIDR,
			"ipv6_gateway":   item.IPv6Gateway,
		}

		subnets = append(subnets, subnet)
		ids = append(ids, item.ID)
	}
	log.Printf("[DEBUG] Subnets List after filter, count=%d :%+v", len(subnets), subnets)

	mErr := d.Set("subnets", subnets)
	if mErr != nil {
		return diag.Errorf("set subnets err:%s", mErr)
	}

	d.SetId(hashcode.Strings(ids))

	return nil
}
