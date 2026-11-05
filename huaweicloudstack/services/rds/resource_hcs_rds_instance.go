package rds

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"

	upstreamRds "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/rds"
)

// @API RDS PUT /v3/{project_id}/instances/{instance_id}/failover/strategy
// @API RDS PUT /v3/{project_id}/instances/{instance_id}/ops-window
func ResourceRdsInstance() *schema.Resource {
	r := upstreamRds.ResourceRdsInstance()

	appendedSwitchStrategy := appendSwitchStrategySchema(r)
	appendedMaintainWindow := appendMaintainWindowSchema(r)

	upstreamUpdate := r.UpdateContext
	r.UpdateContext = func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		diags := upstreamUpdate(ctx, d, meta)
		if len(diags) > 0 && diags.HasError() {
			return diags
		}

		if appendedSwitchStrategy {
			if err := updateRdsInstanceSwitchStrategy(ctx, d, meta, true); err != nil {
				return diag.FromErr(err)
			}
		}
		if appendedMaintainWindow {
			if err := updateRdsInstanceMaintainWindow(ctx, d, meta, true); err != nil {
				return diag.FromErr(err)
			}
		}
		return diags
	}

	upstreamCreate := r.CreateContext
	r.CreateContext = func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		diags := upstreamCreate(ctx, d, meta)
		if len(diags) > 0 && diags.HasError() {
			return diags
		}

		// Skip the request when the strategy is the server-side default (reliability),
		// matching the upstream provider's create behavior and saving an API call.
		if appendedSwitchStrategy {
			if v, ok := d.GetOk("switch_strategy"); ok && v.(string) != "reliability" {
				if err := updateRdsInstanceSwitchStrategy(ctx, d, meta, false); err != nil {
					return diag.FromErr(err)
				}
			}
		}
		if appendedMaintainWindow {
			if err := updateRdsInstanceMaintainWindow(ctx, d, meta, false); err != nil {
				return diag.FromErr(err)
			}
		}
		return diags
	}

	upstreamRead := r.ReadContext
	r.ReadContext = func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		diags := upstreamRead(ctx, d, meta)
		if diags.HasError() || d.Id() == "" {
			return diags
		}
		if !appendedSwitchStrategy && !appendedMaintainWindow {
			return diags
		}

		client, err := getRdsV3Client(meta, d)
		if err != nil {
			return diag.Errorf("error creating RDS client: %s", err)
		}
		instance, err := GetRdsInstanceByID(client, d.Id())
		if err != nil {
			return diag.Errorf("error getting RDS instance: %s", err)
		}

		var mErr *multierror.Error
		if appendedSwitchStrategy {
			if v := utils.PathSearch("switch_strategy", instance, "").(string); v != "" {
				mErr = multierror.Append(mErr, d.Set("switch_strategy", v))
			}
		}
		// The API returns the maintain window as a single string, e.g. "06:00-07:00".
		if appendedMaintainWindow {
			if v := utils.PathSearch("maintenance_window", instance, "").(string); v != "" {
				maintainWindow := strings.Split(v, "-")
				if len(maintainWindow) == 2 {
					mErr = multierror.Append(mErr,
						d.Set("maintain_begin", maintainWindow[0]),
						d.Set("maintain_end", maintainWindow[1]))
				}
			}
		}
		return diag.FromErr(mErr.ErrorOrNil())
	}

	return r
}

func appendSwitchStrategySchema(r *schema.Resource) bool {
	if _, ok := r.Schema["switch_strategy"]; ok {
		return false
	}

	r.Schema["switch_strategy"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	}
	return true
}

// @API RDS PUT /v3/{project_id}/instances/{instance_id}/failover/strategy
func updateRdsInstanceSwitchStrategy(ctx context.Context, d *schema.ResourceData, meta interface{}, waitInstanceReady bool) error {
	if !d.HasChanges("switch_strategy") {
		return nil
	}

	client, err := getRdsV3Client(meta, d)
	if err != nil {
		return fmt.Errorf("error creating RDS client: %s", err)
	}

	_, err = updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:             "v3/{project_id}/instances/{instance_id}/failover/strategy",
		httpMethod:          "PUT",
		pathParams:          map[string]string{"instance_id": d.Id()},
		updateBodyParams:    buildUpdateInstanceSwitchStrategyBodyParams(d),
		isRetry:             true,
		timeout:             schema.TimeoutUpdate,
		isWaitInstanceReady: waitInstanceReady,
	})
	if err != nil {
		return fmt.Errorf("error updating RDS instance (%s) switch strategy: %s", d.Id(), err)
	}
	return nil
}

func buildUpdateInstanceSwitchStrategyBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"repairStrategy": d.Get("switch_strategy"),
	}
}

// @API RDS PUT /v3/{project_id}/instances/{instance_id}/ops-window
func updateRdsInstanceMaintainWindow(ctx context.Context, d *schema.ResourceData, meta interface{}, waitInstanceReady bool) error {
	if !d.HasChanges("maintain_begin", "maintain_end") {
		return nil
	}

	client, err := getRdsV3Client(meta, d)
	if err != nil {
		return fmt.Errorf("error creating RDS client: %s", err)
	}

	_, err = updateRdsInstanceField(ctx, d, client, updateInstanceFieldParams{
		httpUrl:             "v3/{project_id}/instances/{instance_id}/ops-window",
		httpMethod:          "PUT",
		pathParams:          map[string]string{"instance_id": d.Id()},
		updateBodyParams:    buildUpdateInstanceMaintainWindowBodyParams(d),
		isWaitInstanceReady: waitInstanceReady,
	})
	if err != nil {
		return fmt.Errorf("error updating RDS instance (%s) maintain window: %s", d.Id(), err)
	}
	return nil
}

func buildUpdateInstanceMaintainWindowBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"start_time": d.Get("maintain_begin").(string),
		"end_time":   d.Get("maintain_end").(string),
	}
}

func getRdsV3Client(meta interface{}, d *schema.ResourceData) (*golangsdk.ServiceClient, error) {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	return cfg.NewServiceClient("rds", region)
}
