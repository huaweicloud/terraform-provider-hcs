package rds

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

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
