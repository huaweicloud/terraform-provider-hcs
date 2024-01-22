package evs

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/evs/v2/snapshots"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
)

func DataSourceEvsSnapshots() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceSnapshotsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapshots": {
				Type:     schema.TypeList,
				Elem:     snapshotSchema(),
				Computed: true,
			},
		},
	}
}

func snapshotSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
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
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	return &sc
}

func datasourceSnapshotsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	client, err := cfg.BlockStorageV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloudStack EVS v2 client: %s", err)
	}

	pages, err := snapshots.List(client, buildEVSSnapshotsQueryParams(d, cfg)).AllPages()
	if err != nil {
		return fmtp.DiagErrorf("an error occurred while fetching the pages of the EVS snapshots: %s", err)
	}
	sps, err := snapshots.ExtractSnapshots(pages)
	if err != nil {
		return fmtp.DiagErrorf("error getting the EVS snapshot list form server: %s", err)
	}

	sMap, ids, err := sourceEvsSnapshots(sps)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(hashcode.Strings(ids))
	if err = d.Set("snapshots", sMap); err != nil {
		return fmtp.DiagErrorf("error saving the detailed information of the EVS snapshots to state: %s", err)
	}
	return nil
}

func buildEVSSnapshotsQueryParams(d *schema.ResourceData, cfg *config.HcsConfig) snapshots.ListOpts {
	result := snapshots.ListOpts{
		AvailabilityZone: d.Get("availability_zone").(string),
		Status:           d.Get("status").(string),
		ID:               d.Get("snapshot_id").(string),
		Name:             d.Get("name").(string),
		VolumeID:         d.Get("volume_id").(string),
	}
	enterpriseProjectID := cfg.DataGetEnterpriseProjectID(d)
	if enterpriseProjectID != "all_granted_eps" {
		result.EnterpriseProjectIDs = "['" + enterpriseProjectID + "']"
	}
	return result
}

func sourceEvsSnapshots(sps []snapshots.Snapshot) ([]map[string]interface{}, []string, error) {
	result := make([]map[string]interface{}, len(sps))
	ids := make([]string, len(sps))

	for i, snapshot := range sps {
		sMap := map[string]interface{}{
			"id":          snapshot.ID,
			"description": snapshot.Description,
			"name":        snapshot.Name,
			"size":        snapshot.Size,
			"status":      snapshot.Status,
			"created_at":  snapshot.CreatedAt.Format(time.RFC3339),
			"updated_at":  snapshot.UpdatedAt.Format(time.RFC3339),
			"metadata":    snapshot.Metadata,
			"volume_id":   snapshot.VolumeID,
		}
		result[i] = sMap
		ids[i] = snapshot.ID
	}
	return result, ids, nil
}
