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
//
// ResourceRdsInstance wraps the upstream hcs_rds_instance resource to add three arguments that
// are missing from the upstream version locked by go.mod (commit 51169dadfcdb, 2023-12-13):
// switch_strategy, maintain_begin and maintain_end.
//
// The Terraform Plugin SDK v2 cannot append schema fields to a Resource returned by another
// package, so instead of copying the 1500+ line upstream implementation, this function:
//  1. gets the upstream Resource via upstreamRds.ResourceRdsInstance();
//  2. appends the three schema fields;
//  3. wraps CreateContext/UpdateContext to apply them through separate PUT requests.
//
// The three arguments are not part of the create-instance API body; they are updated afterwards
// via dedicated PUT requests, which is the same approach the upstream provider takes.
func ResourceRdsInstance() *schema.Resource {
	r := upstreamRds.ResourceRdsInstance()

	// Append the three new schema fields that are missing from the locked upstream version.
	// These fields are ported from the public cloud provider's current version.
	r.Schema["switch_strategy"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
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

	// Wrap the upstream UpdateContext so the new arguments are applied after the upstream update
	// completes. waitInstanceReady is true here because the upstream UpdateContext may return
	// while the instance is still in UPDATING state.
	upstreamUpdate := r.UpdateContext
	r.UpdateContext = func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		diags := upstreamUpdate(ctx, d, meta)
		if len(diags) > 0 && diags.HasError() {
			return diags
		}

		if err := updateRdsInstanceSwitchStrategy(ctx, d, meta, true); err != nil {
			return diag.FromErr(err)
		}
		if err := updateRdsInstanceMaintainWindow(ctx, d, meta, true); err != nil {
			return diag.FromErr(err)
		}
		return diags
	}

	// Wrap the upstream CreateContext so the new arguments are applied right after the instance is
	// created. waitInstanceReady is false here because the upstream create already waits for ACTIVE.
	upstreamCreate := r.CreateContext
	r.CreateContext = func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		diags := upstreamCreate(ctx, d, meta)
		if len(diags) > 0 && diags.HasError() {
			return diags
		}

		if err := updateRdsInstanceSwitchStrategy(ctx, d, meta, false); err != nil {
			return diag.FromErr(err)
		}
		if err := updateRdsInstanceMaintainWindow(ctx, d, meta, false); err != nil {
			return diag.FromErr(err)
		}
		return diags
	}

	// Wrap the upstream ReadContext to backfill the three new fields.
	//
	// The locked upstream version (commit 51169dadfcdb, 2023-12-13) never sets these fields in
	// Read, which causes state drift: when an update fails, Terraform has already written the
	// desired value into state, while the real instance keeps its previous value. Backfilling
	// here makes the next refresh correct the state, matching the upstream provider's behavior.
	upstreamRead := r.ReadContext
	r.ReadContext = func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		diags := upstreamRead(ctx, d, meta)
		if diags.HasError() || d.Id() == "" {
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
		if v := utils.PathSearch("switch_strategy", instance, "").(string); v != "" {
			mErr = multierror.Append(mErr, d.Set("switch_strategy", v))
		}
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
