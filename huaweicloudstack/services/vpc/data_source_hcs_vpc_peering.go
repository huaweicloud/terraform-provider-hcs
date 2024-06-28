package vpc

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	v1peerings "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/peerings"
)

func DataSourceVpcPeering() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcPeeringRead,

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
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peer_vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVpcPeeringRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	region := hcsConfig.GetRegion(d)
	vpcPeeringClient, err := hcsConfig.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC Peering client: %s", err)
	}

	listOpts := v1peerings.ListOpts{
		ID:   d.Get("id").(string),
		Name: d.Get("name").(string),
	}

	n, err := v1peerings.List(vpcPeeringClient, listOpts).ExtractList()
	if len(n) <= 0 || err != nil {
		return diag.Errorf("error querying VPC Peering: %s", err)
	}
	d.SetId(n[0].ID)
	mErr := multierror.Append(nil,
		d.Set("name", n[0].Name),
		d.Set("vpc_id", n[0].RequesterVpcInfo.VpcId),
		d.Set("peer_vpc_id", n[0].AccepterVpcInfo.VpcId),
		d.Set("region", hcsConfig.GetRegion(d)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
