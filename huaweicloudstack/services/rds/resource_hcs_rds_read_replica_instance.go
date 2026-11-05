package rds

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"

	upstreamRds "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rds"
)

// @API RDS PUT /v3/{project_id}/instances/{instance_id}/ops-window
func ResourceRdsReadReplicaInstance() *schema.Resource {
	r := upstreamRds.ResourceRdsReadReplicaInstance()

	appended := appendMaintainWindowSchema(r)

	upstreamCreate := r.CreateContext
	r.CreateContext = func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		diags := upstreamCreate(ctx, d, meta)
		if len(diags) > 0 && diags.HasError() {
			return diags
		}
		if !appended {
			return diags
		}

		if err := updateRdsInstanceMaintainWindow(ctx, d, meta, false); err != nil {
			return diag.FromErr(err)
		}
		return diags
	}

	upstreamUpdate := r.UpdateContext
	r.UpdateContext = func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		diags := upstreamUpdate(ctx, d, meta)
		if len(diags) > 0 && diags.HasError() {
			return diags
		}
		if !appended {
			return diags
		}

		if err := updateRdsInstanceMaintainWindow(ctx, d, meta, true); err != nil {
			return diag.FromErr(err)
		}
		return diags
	}

	upstreamRead := r.ReadContext
	r.ReadContext = func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		diags := upstreamRead(ctx, d, meta)
		if diags.HasError() || d.Id() == "" || !appended {
			return diags
		}

		client, err := getRdsV3Client(meta, d)
		if err != nil {
			return diag.Errorf("error creating RDS client: %s", err)
		}
		instance, err := GetRdsInstanceByID(client, d.Id())
		if err != nil {
			return diag.Errorf("error getting RDS read replica instance: %s", err)
		}

		var mErr *multierror.Error
		// The API returns the maintain window as a single string, e.g. "06:00-07:00".
		if v := utils.PathSearch("maintenance_window", instance, "").(string); v != "" {
			maintainWindow := strings.Split(v, "-")
			if len(maintainWindow) == 2 {
				mErr = multierror.Append(mErr,
					d.Set("maintain_begin", maintainWindow[0]),
					d.Set("maintain_end", maintainWindow[1]))
			}
		}
		return diag.FromErr(mErr.ErrorOrNil())
	}

	return r
}

func appendMaintainWindowSchema(r *schema.Resource) bool {
	if _, ok := r.Schema["maintain_begin"]; ok {
		return false
	}
	if _, ok := r.Schema["maintain_end"]; ok {
		return false
	}

	r.Schema["maintain_begin"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	}
	r.Schema["maintain_end"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		RequiredWith: []string{"maintain_begin"},
	}
	return true
}
