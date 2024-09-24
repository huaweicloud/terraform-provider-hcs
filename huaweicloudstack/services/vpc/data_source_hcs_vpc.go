package vpc

import (
	"context"
	v3Vpcs "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v3/vpcs"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/vpcs"
)

func DataSourceVpcV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcV1Read,

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
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"secondary_cidrs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func dataSourceVpcV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	vpcClient, err := hcsConfig.NetworkingV1Client(hcsConfig.GetRegion(d))
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

	refinedVpcs, err := vpcs.List(vpcClient, listOpts)
	if err != nil {
		return diag.Errorf("unable to retrieve vpcs: %s", err)
	}

	if len(refinedVpcs) < 1 {
		return diag.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedVpcs) > 1 {
		return diag.Errorf("your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Vpc := refinedVpcs[0]

	log.Printf("[INFO] Retrieved Vpc using given filter %s: %+v", Vpc.ID, Vpc)
	d.SetId(Vpc.ID)

	d.Set("region", hcsConfig.GetRegion(d))
	d.Set("name", Vpc.Name)
	d.Set("enterprise_project_id", Vpc.EnterpriseProjectID)
	d.Set("cidr", Vpc.CIDR)
	d.Set("status", Vpc.Status)

	var s []map[string]interface{}
	for _, route := range Vpc.Routes {
		mapping := map[string]interface{}{
			"destination": route.DestinationCIDR,
			"nexthop":     route.NextHop,
		}
		s = append(s, mapping)
	}
	d.Set("routes", s)
	vpcV3Client, v3Err := hcsConfig.NetworkingV3Client(hcsConfig.GetRegion(d))
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
