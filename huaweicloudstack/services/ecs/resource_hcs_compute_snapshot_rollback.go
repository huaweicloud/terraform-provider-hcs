package ecs

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/ecs/v1/snapshots"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/logp"
)

func ResourceComputeSnapshotRollback() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeSnapshotRollbackCreate,
		ReadContext:   resourceComputeSnapshotRollbackRead,
		DeleteContext: resourceComputeSnapshotRollbackDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"snapshot_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeSnapshotRollbackCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	ecsV2Client, err := cfg.ComputeV2Client(region)
	if err != nil {
		return diag.Errorf("error creating compute V2 client: %s", err)
	}

	snapshotId := d.Get("snapshot_id").(string)

	rollbackOpts := snapshots.RollBackInstanceSnapshotOpts{
		ImageRef: snapshotId,
	}

	logp.Printf("[DEBUG] snapshot rollback options: %#v", rollbackOpts)
	serverId := d.Get("instance_id").(string)

	imsClient, err := cfg.ImageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ims V2 client: %s", err)
	}

	// snapshot is existed ?
	queryImage, err := snapshots.Get(imsClient, serverId, snapshotId)
	if queryImage.Id == "" {
		return diag.Errorf("failed to query snapshot: %s", err)
	}

	jobStatus, err := snapshots.Rollback(ecsV2Client, serverId, rollbackOpts).ExtractJobStatus()
	if err != nil {
		return diag.Errorf("failed to rollback an snapshot (%s) for instance (%s). err: %s",
			snapshotId, serverId, err)
	}

	ecsV1Client, err := cfg.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute V1 client: %s", err)
	}

	// The value same as the timeout interval for creating an instance.
	timeout := 30 * time.Minute
	if err := snapshots.WaitForJobSuccess(ecsV1Client, int(timeout/time.Second), jobStatus.JobID); err != nil {
		return diag.Errorf("failed to wait snapshot (%s) to be rollbacked for the instance (%s): %s",
			snapshotId, serverId, err)
	}

	d.SetId(serverId + "." + snapshotId)
	return resourceComputeSnapshotRollbackRead(ctx, d, meta)
}

func resourceComputeSnapshotRollbackRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := d.Set("snapshot_id", "")
	if err != nil {
		return diag.Errorf("failed to set snapshot_id for the instance (%s): %s",
			d.Get("instance_id").(string), err)
	}
	return nil
}

func resourceComputeSnapshotRollbackDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
