package ecs

import (
	"context"
	"fmt"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/ecs/v1/cloudservers"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
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
			StateContext: resourceComputeSnapshotImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
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
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	imageV2Client, err := cfg.ImageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating compute V2 client: %s", err)
	}
	var instanceId string
	var snapshotId string
	if strings.Contains(d.Id(), "/") {
		parts := strings.SplitN(d.Id(), "/", 2)
		instanceId = parts[0]
		snapshotId = parts[1]
	} else {
		instanceId = d.Get("instance_id").(string)
		snapshotId = d.Id()
	}
	snapshot, err := snapshots.Get(imageV2Client, instanceId, snapshotId)
	if err != nil {
		return diag.Errorf("error query snapshot: %s", err)
	} else if snapshot.Id == "" {
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] Retrieved Snapshot %s: %#v", d.Id(), snapshot)
	d.Set("instance_id", d.Get("instance_id").(string))
	d.Set("name", snapshot.Name)
	return nil
}

func resourceComputeSnapshotDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	ecsClient, err := cfg.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}
	server, err := cloudservers.Get(ecsClient, d.Get("instance_id").(string)).Extract()
	var images []string
	images = append(images, d.Id())
	deleteOpts := snapshots.DeleteOpts{
		Images:          images,
		AvailableZone:   server.AvailabilityZone,
		Region:          region,
		IsSnapShotImage: "true",
	}
	imageV2Client, err := cfg.ImageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating image V2 client: %s", err)
	}
	n, err := snapshots.Delete(imageV2Client, deleteOpts).ExtractJobResponse()
	if err != nil {
		return diag.Errorf("error deleting snapshot: %s", err)
	}
	ecsV1Client, err := cfg.ComputeV1Client(region)
	if err != nil {
		return diag.Errorf("error creating compute V1 client: %s", err)
	}
	if err := snapshots.WaitForJobSuccess(ecsV1Client, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func resourceComputeSnapshotImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for compute snapshot, must be <instance_id>/<snapshot_id>")
	}
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	imageV2Client, err := cfg.ImageV2Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating compute client: %s", err)
	}
	queryImage, err := snapshots.Get(imageV2Client, parts[0], parts[1])
	if queryImage.Id == "" {
		return nil, common.CheckDeleted(d, err, "compute snapshot")
	}
	d.Set("instance_id", queryImage.SnapshotFromInstance)
	d.Set("snapshot_id", queryImage.Id)
	d.Set("name", queryImage.Name)
	return []*schema.ResourceData{d}, nil
}
