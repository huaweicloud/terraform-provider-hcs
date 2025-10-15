package vpc

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v2/dc/direct_connects"
	"log"
)

func ResourceDirectConnect() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDirectConnectCreate,
		ReadContext:   resourceDirectConnectRead,
		UpdateContext: resourceDirectConnectUpdate,
		DeleteContext: resourceDirectConnectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hosting_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dc_provider": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"hosted"}, false),
			},
			"peer_location": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"group": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDirectConnectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	dcClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Direct Connect client: %s", err)
	}

	createOpts := direct_connects.CreateOpts{
		Name:         d.Get("name").(string),
		Type:         d.Get("type").(string),
		HostingId:    d.Get("hosting_id").(string),
		PeerLocation: d.Get("peer_location").(string),
		Description:  d.Get("description").(string),
	}

	n, err := direct_connects.Create(dcClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating Direct Connect: %s", err)
	}

	d.SetId(n.ID)
	log.Printf("[DEBUG] Virtual Interface ID: %s", n.ID)
	return resourceDirectConnectRead(ctx, d, meta)
}

func resourceDirectConnectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	dcClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Direct Connect client: %s", err)
	}

	n, err := direct_connects.Get(dcClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Direct Connect Connection")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", n.Name),
		d.Set("status", n.Status),
		d.Set("description", n.Description),
		d.Set("hosting_id", n.HostingId),
		d.Set("dc_provider", n.Provider),
		d.Set("type", n.Type),
		d.Set("peer_location", n.PeerLocation),
		d.Set("group", n.Group),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Direct Connect fields: %s", err)
	}

	return nil
}

func resourceDirectConnectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	dcClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Direct Connect client: %s", err)
	}

	updateOpts := direct_connects.UpdateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	_, err = direct_connects.Update(dcClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating Direct Connect: %s", err)
	}

	return resourceDirectConnectRead(ctx, d, meta)
}

func resourceDirectConnectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	dcClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Direct Connect client: %s", err)
	}

	err = direct_connects.Delete(dcClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting Direct Connect: %s", err)
	}

	return nil
}
