package evs

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/blockstorage/extensions/volumeactions"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/blockstorage/v2/volumes"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/common/tags"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
)

func ResourceEvsVolume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvsVolumeCreate,
		ReadContext:   resourceEvsVolumeRead,
		UpdateContext: resourceEvsVolumeUpdate,
		DeleteContext: resourceEvsVolumeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_vol_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"multiattach": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"bootable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"attachments": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attached_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attached_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: resourceVolumeV2AttachmentHash,
			},
			"wwn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceContainerMetadataV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("metadata").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func resourceEvsVolumeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	blockStorageClient, err := cfg.BlockStorageV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloudStack block storage client: %s", err)
	}

	createOpts := &volumes.CreateOpts{
		AvailabilityZone:    d.Get("availability_zone").(string),
		Description:         d.Get("description").(string),
		ImageID:             d.Get("image_id").(string),
		Metadata:            resourceContainerMetadataV2(d),
		Name:                d.Get("name").(string),
		Size:                d.Get("size").(int),
		SnapshotID:          d.Get("snapshot_id").(string),
		SourceVolID:         d.Get("source_vol_id").(string),
		VolumeType:          d.Get("volume_type").(string),
		Multiattach:         d.Get("multiattach").(bool),
		EnterpriseProjectID: cfg.GetEnterpriseProjectID(d),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := volumes.Create(blockStorageClient, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloudStack volume: %s", err)
	}
	log.Printf("[INFO] Volume ID: %s", v.ID)

	// Wait for the volume to become available.
	log.Printf(
		"[DEBUG] Waiting for volume (%s) to become available",
		v.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"downloading", "creating"},
		Target:     []string{"available"},
		Refresh:    VolumeV2StateRefreshFunc(blockStorageClient, v.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf(
			"Error waiting for volume (%s) to become ready: %s",
			v.ID, err)
	}

	// Store the ID now
	d.SetId(v.ID)

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		tagList := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(blockStorageClient, "cloudvolumes", v.ID, tagList).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of evs %s: %s", v.ID, tagErr)
		}
	}

	return resourceEvsVolumeRead(ctx, d, meta)
}

func resourceEvsVolumeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	blockStorageClient, err := cfg.BlockStorageV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloudStack block storage client: %s", err)
	}

	v, err := volumes.Get(blockStorageClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "volume")
	}

	log.Printf("[DEBUG] Retrieved volume %s: %+v", d.Id(), v)

	d.Set("size", v.Size)
	d.Set("description", v.Description)
	d.Set("availability_zone", v.AvailabilityZone)
	d.Set("name", v.Name)
	d.Set("snapshot_id", v.SnapshotID)
	d.Set("source_vol_id", v.SourceVolID)
	d.Set("volume_type", v.VolumeType)
	d.Set("wwn", v.WWN)
	d.Set("status", v.Status)
	d.Set("created_at", v.CreatedAt)
	d.Set("updated_at", v.UpdatedAt)
	d.Set("metadata", v.Metadata)
	d.Set("multiattach", v.Multiattach)
	d.Set("tags", v.Tags)
	d.Set("enterprise_project_id", v.EnterpriseProjectID)
	d.Set("region", cfg.GetRegion(d))

	attachments := make([]map[string]interface{}, len(v.Attachments))
	for i, attachment := range v.Attachments {
		attachments[i] = make(map[string]interface{})
		attachments[i]["id"] = attachment.AttachmentID
		attachments[i]["server_id"] = attachment.ServerID
		attachments[i]["device_name"] = attachment.Device
		attachments[i]["attached_at"] = attachment.AttachedAt
		attachments[i]["attached_mode"] = v.Metadata["attached_mode"]
		log.Printf("[DEBUG] attachment: %v", attachment)
	}
	d.Set("attachments", attachments)
	bootable, err := strconv.ParseBool(v.Bootable)
	if err != nil {
		return fmtp.DiagErrorf("The bootable of volume (%s) connot be converted to bool.",
			v.ID)
	}
	d.Set("bootable", bootable)
	return nil
}

func resourceEvsVolumeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	blockStorageClient, err := cfg.BlockStorageV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloudStack block storage client: %s", err)
	}

	updateOpts := volumes.UpdateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	if d.HasChange("metadata") {
		updateOpts.Metadata = resourceVolumeMetadataV2(d)
	}

	if d.HasChange("size") {
		extendOpts := volumeactions.ExtendSizeOpts{
			NewSize: d.Get("size").(int),
		}

		err = volumeactions.ExtendSize(blockStorageClient, d.Id(), extendOpts).ExtractErr()
		if err != nil {
			return fmtp.DiagErrorf("Error extending huaweicloudstack_blockstorage_volume_v2 %s size: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"extending"},
			Target:     []string{"available", "in-use"},
			Refresh:    VolumeV2StateRefreshFunc(blockStorageClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err := stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmtp.DiagErrorf(
				"Error waiting for huaweicloudstack_blockstorage_volume_v2 %s to become ready: %s", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(blockStorageClient, d, "cloudvolumes", d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of volume:%s, err:%s", d.Id(), tagErr)
		}
	}

	_, err = volumes.Update(blockStorageClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error updating HuaweiCloudStack volume: %s", err)
	}

	return resourceEvsVolumeRead(ctx, d, meta)
}

func resourceEvsVolumeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	blockStorageClient, err := cfg.BlockStorageV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloudStack block storage client: %s", err)
	}

	v, err := volumes.Get(blockStorageClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "volume")
	}

	// make sure this volume is detached from all instances before deleting
	if len(v.Attachments) > 0 {
		log.Printf("[DEBUG] detaching volumes")
		if computeClient, err := cfg.ComputeV2Client(cfg.GetRegion(d)); err != nil {
			return fmtp.DiagErrorf(
				"Error get volume (%s) attachments, errMsg: %s",
				d.Id(), err)
		} else {
			for _, volumeAttachment := range v.Attachments {
				log.Printf("[DEBUG] Attachment: %v", volumeAttachment)
				err := volumeattach.Delete(computeClient, volumeAttachment.ServerID, volumeAttachment.ID).ExtractErr()
				if err != nil {
					return fmtp.DiagErrorf(
						"Error detach for volume (%s), errMsg: %s",
						d.Id(), err)
				}
			}

			stateConf := &resource.StateChangeConf{
				Pending:    []string{"in-use", "attaching", "detaching"},
				Target:     []string{"available"},
				Refresh:    VolumeV2StateRefreshFunc(blockStorageClient, d.Id()),
				Timeout:    d.Timeout(schema.TimeoutDelete),
				Delay:      10 * time.Second,
				MinTimeout: 3 * time.Second,
			}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return fmtp.DiagErrorf(
					"Error waiting for volume (%s) to become available: %s",
					d.Id(), err)
			}
		}
	}

	// It's possible that this volume was used as a boot device and is currently
	// in a "deleting" state from when the instance was terminated.
	// If this is true, just move on. It'll eventually delete.
	if v.Status != "deleting" {
		if err := volumes.Delete(blockStorageClient, d.Id(), nil).ExtractErr(); err != nil {
			return common.CheckDeletedDiag(d, err, "volume")
		}
	}

	// Wait for the volume to delete before moving on.
	log.Printf("[DEBUG] Waiting for volume (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting", "downloading", "available"},
		Target:     []string{"deleted"},
		Refresh:    VolumeV2StateRefreshFunc(blockStorageClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf(
			"Error waiting for volume (%s) to delete: %s",
			d.Id(), err)
	}

	d.SetId("")
	return nil
}

func resourceVolumeMetadataV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("metadata").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

// VolumeV2StateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an HuaweiCloudStack volume.
func VolumeV2StateRefreshFunc(client *golangsdk.ServiceClient, volumeID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := volumes.Get(client, volumeID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "deleted", nil
			}
			return nil, "", err
		}

		if v.Status == "error" {
			return v, v.Status, fmt.Errorf("There was an error creating the volume. " +
				"Please check with your cloud admin or check the Block Storage " +
				"API logs to see why this error occurred.")
		}

		return v, v.Status, nil
	}
}

func resourceVolumeV2AttachmentHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if m["instance_id"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["instance_id"].(string)))
	}
	return hashcode.String(buf.String())
}
