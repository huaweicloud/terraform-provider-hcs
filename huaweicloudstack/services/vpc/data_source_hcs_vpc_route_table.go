package vpc

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/routetables"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
)

// @API VPC GET /v1/{project_id}/routetables/{id}
// @API VPC GET /v1/{project_id}/routetables
func DataSourceVPCRouteTable() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcRouteTableRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"route": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nexthop": {
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

func dataSourceVpcRouteTableRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	region := hcsConfig.GetRegion(d)
	vpcClient, err := hcsConfig.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	listOpts := routetables.ListOpts{
		VpcID: d.Get("vpc_id").(string),
		ID:    d.Get("id").(string),
	}
	pages, err := routetables.List(vpcClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("unable to retrieve route tables: %s", err)
	}

	allRouteTables, err := routetables.ExtractRouteTables(pages)
	if err != nil {
		return diag.Errorf("unable to extract route tables: %s", err)
	}

	if len(allRouteTables) < 1 {
		return diag.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	var rtbID string
	if v, ok := d.GetOk("name"); ok {
		filterName := v.(string)
		for _, rtb := range allRouteTables {
			if filterName == rtb.Name {
				rtbID = rtb.ID
				break
			}
		}
	} else {
		// find the default route table if name was not specified
		for _, rtb := range allRouteTables {
			if rtb.Default {
				rtbID = rtb.ID
				break
			}
		}
	}

	if rtbID == "" {
		return diag.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	// call Get API to retrieve more details about the route table
	routeTable, err := routetables.Get(vpcClient, rtbID).Extract()
	if err != nil {
		return diag.Errorf("unable to retrieve route table %s: %s", rtbID, err)
	}

	log.Printf("[DEBUG] Retrieved VPC route table %s: %+v", rtbID, routeTable)
	d.SetId(rtbID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("vpc_id", routeTable.VpcID),
		d.Set("name", routeTable.Name),
		d.Set("description", routeTable.Description),
		d.Set("default", routeTable.Default),
		d.Set("subnets", expandVpcRTSubnets(routeTable.Subnets)),
		d.Set("route", expandVpcRTRoutes(routeTable.Routes)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving VPC route table: %s", err)
	}

	return nil
}
