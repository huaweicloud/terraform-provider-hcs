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

func ResourceComputeSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeSnapshotCreate,
		ReadContext:   resourceComputeSnapshotRead,
		DeleteContext: resourceComputeSnapshotDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	ecsV2Client, err := cfg.ComputeV2Client(region)
	if err != nil {
		return diag.Errorf("error creating compute V21 client: %s", err)
	}

	createOpts := snapshots.CreateInstanceSnapshotOpts{
		Name:             d.Get("name").(string),
		InstanceSnapshot: "true",
		ServerId:         d.Get("instance_id").(string),
	}

	logp.Printf("[DEBUG] create instance snapshot options: %#v", createOpts)

	jobStatus, err := snapshots.Create(ecsV2Client, createOpts).ExtractJobStatus()
	if err != nil {
		return diag.Errorf("failed to create an snapshot for instance (%s). err: %s", createOpts.ServerId, err)
	}

	ecsV1Client, err := cfg.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute V1 client: %s", err)
	}

	// The duration of an ECS snapshot depends on the volume capacity, used capacity, and storage performance.
	// The task timeout interval of the cloud service interface is 13.88 hours.
	timeout := 14 * time.Hour
	if err := snapshots.WaitForJobSuccess(ecsV1Client, int(timeout/time.Second), jobStatus.JobID); err != nil {
		return diag.Errorf("failed to wait for the instance (%s) snapshot to be created: %s", createOpts.ServerId, err)
	}

	imageId, err := snapshots.GetJobEntity(ecsV1Client, jobStatus.JobID, "image_id")
	if err != nil {
		return diag.Errorf("failed to get snapshot id for the instance (%s): %s", createOpts.ServerId, err)
	}

	d.SetId(imageId.(string))
	return resourceComputeSnapshotRead(ctx, d, meta)
}

func resourceComputeSnapshotRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := d.Set("name", "")
	if err != nil {
		return diag.Errorf("failed to set snapshot name for instance (%s): %s", d.Get("instance_id").(string), err)
	}
	return nil
}

func resourceComputeSnapshotDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
