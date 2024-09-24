package vpc

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	v1peeringroute "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/peeringsroute"
)

func ResourceVpcPeeringRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcPeeringRouteCreate,
		ReadContext:   resourceVpcPeeringRouteRead,
		UpdateContext: resourceVpcPeeringRouteUpdate,
		DeleteContext: resourceVpcPeeringRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{ // request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"peering_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 200,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination": {
							Type:     schema.TypeString,
							Required: true,
						},
						"nexthop": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceVpcPeeringRouteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	region := hcsConfig.GetRegion(d)
	vpcId := d.Get("vpc_id").(string)
	vpcPeeringClient, err := hcsConfig.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC Peering client: %s", err)
	}

	createOpts := v1peeringroute.CreateOpts{
		Route: buildVpcPeeringRoutes(d),
	}

	_, err = v1peeringroute.Create(vpcPeeringClient, createOpts, vpcId).Extract()
	if err != nil {
		return diag.Errorf("error creating VPC Peering: %s", err)
	}

	d.SetId(d.Get("peering_id").(string) + "/" + d.Get("vpc_id").(string))

	return resourceVpcPeeringRouteRead(ctx, d, meta)
}

func buildVpcPeeringRoutes(d *schema.ResourceData) []v1peeringroute.Route {
	rawRoutes := d.Get("route").(*schema.Set).List()
	routeOpts := make([]v1peeringroute.Route, len(rawRoutes))

	for i, raw := range rawRoutes {
		opts := raw.(map[string]interface{})
		routeOpts[i] = v1peeringroute.Route{
			NextHop:     opts["nexthop"].(string),
			Destination: opts["destination"].(string),
		}
	}

	return routeOpts
}

func resourceVpcPeeringRouteRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	region := hcsConfig.GetRegion(d)
	vpcPeeringClient, err := hcsConfig.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC Peering client: %s", err)
	}

	id := strings.Split(d.Id(), "/")
	if len(id) != 2 {
		return diag.Errorf("ID format error: %s", d.Id())
	}
	listOps := v1peeringroute.ListOpts{
		PeeringId: id[0],
	}

	n, err := v1peeringroute.List(vpcPeeringClient, listOps, id[1]).Extract()
	if err != nil {
		return diag.Errorf("error creating VPC Peering client: %s", err)
	}
	mErr := multierror.Append(nil,
		d.Set("peering_id", id[0]),
		d.Set("vpc_id", id[1]),
		d.Set("route", expandVpcPeeringRoutes(n)),
		d.Set("region", hcsConfig.GetRegion(d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func expandVpcPeeringRoutes(routes []v1peeringroute.Route) []map[string]interface{} {
	rtRules := make([]map[string]interface{}, 0, len(routes))

	for _, item := range routes {
		acessRule := map[string]interface{}{
			"destination": item.Destination,
			"nexthop":     item.NextHop,
		}
		rtRules = append(rtRules, acessRule)
	}

	return rtRules
}

func resourceVpcPeeringRouteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	region := hcsConfig.GetRegion(d)
	id := strings.Split(d.Id(), "/")
	if len(id) != 2 {
		return diag.Errorf("ID format error: %s", d.Id())
	}
	vpcId := id[1]
	vpcPeeringClient, err := hcsConfig.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC Peering client: %s", err)
	}

	routes := buildVpcPeeringRoutes(d)
	if len(routes) == 0 {
		d.SetId("")
		return nil
	}
	createOpts := v1peeringroute.CreateOpts{
		Route: routes,
	}

	_, err = v1peeringroute.Delete(vpcPeeringClient, createOpts, vpcId).Extract()
	if _, ok := err.(golangsdk.ErrDefault400); ok {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("error deleting VPC Peering %s: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func resourceVpcPeeringRouteUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	vpcClient, err := hcsConfig.NetworkingV1Client(hcsConfig.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}
	if d.HasChange("route") {
		oldRoute, newRoute := d.GetChange("route")
		addRaws := newRoute.(*schema.Set).Difference(oldRoute.(*schema.Set))
		delRaws := oldRoute.(*schema.Set).Difference(newRoute.(*schema.Set))

		id := strings.Split(d.Id(), "/")
		if len(id) != 2 {
			return diag.Errorf("ID format error: %s", d.Id())
		}
		vpcId := id[1]

		if delLen := delRaws.Len(); delLen > 0 {
			delRouteOpts := make([]v1peeringroute.Route, delLen)
			for i, item := range delRaws.List() {
				opts := item.(map[string]interface{})
				delRouteOpts[i] = v1peeringroute.Route{
					NextHop:     opts["nexthop"].(string),
					Destination: opts["destination"].(string),
				}
			}
			createOpts := v1peeringroute.CreateOpts{
				Route: delRouteOpts,
			}
			if _, err := v1peeringroute.Delete(vpcClient, createOpts, vpcId).Extract(); err != nil {
				return diag.Errorf("error deleting VPC peering routes: %s", err)
			}
		}

		if addLen := addRaws.Len(); addLen > 0 {
			addRouteOpts := make([]v1peeringroute.Route, addLen)
			for i, item := range addRaws.List() {
				opts := item.(map[string]interface{})
				addRouteOpts[i] = v1peeringroute.Route{
					NextHop:     opts["nexthop"].(string),
					Destination: opts["destination"].(string),
				}
			}
			createOpts := v1peeringroute.CreateOpts{
				Route: addRouteOpts,
			}
			if _, err := v1peeringroute.Create(vpcClient, createOpts, vpcId).Extract(); err != nil {
				return diag.Errorf("error adding VPC peering routes: %s", err)
			}
		}
	}
	return resourceVpcPeeringRouteRead(ctx, d, meta)
}
