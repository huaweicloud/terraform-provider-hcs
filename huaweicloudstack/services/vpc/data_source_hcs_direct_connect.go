package vpc

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v2/dc/direct_connects"
	"log"
)

func DataSourceDirectConnect() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataResourceDirectConnectRead,

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
			"hosting_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dc_provider": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peer_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataResourceDirectConnectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	dcClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Direct Connect client: %s", err)
	}

	listOpts := direct_connects.ListOpts{
		ID:       d.Get("id").(string),
		Name:     d.Get("name").(string),
		Provider: d.Get("dc_provider").(string),
		Type:     "hosted",
	}
	n, err := direct_connects.List(dcClient, listOpts)
	if len(n) <= 0 || err != nil {
		return diag.Errorf("error querying Direct Connect: %s", err)
	}
	log.Printf("[INFO] Retrieved Direct Connect using given filter %s: %+v", n[0].ID, n[0])

	d.SetId(n[0].ID)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", n[0].Name),
		d.Set("status", n[0].Status),
		d.Set("description", n[0].Description),
		d.Set("hosting_id", n[0].HostingId),
		d.Set("dc_provider", n[0].Provider),
		d.Set("type", n[0].Type),
		d.Set("peer_location", n[0].PeerLocation),
		d.Set("group", n[0].Group),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Direct Connect fields: %s", err)
	}

	return nil
}
