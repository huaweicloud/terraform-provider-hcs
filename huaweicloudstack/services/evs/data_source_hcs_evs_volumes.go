package evs

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/blockstorage/v2/volumes"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
)

func DataSourceEvsVolumesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEvsVolumesV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"multiattach": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attachments": {
							Type:     schema.TypeList,
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
									"device_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"server_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bootable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"multiattach": {
							Type:     schema.TypeBool,
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
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metadata": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"wwn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildQueryOpts(d *schema.ResourceData, cfg *config.HcsConfig) volumes.ListOpts {
	result := volumes.ListOpts{
		AvailabilityZone:    d.Get("availability_zone").(string),
		EnterpriseProjectID: cfg.DataGetEnterpriseProjectID(d),
		Status:              d.Get("status").(string),
		Name:                d.Get("name").(string),
		MetadataOrigin:      d.Get("metadata").(map[string]interface{}),
		Tags:                resourceContainerTags(d),
	}
	return result
}

func resourceContainerTags(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("tags").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func sourceEvsAttachment(attachements []volumes.Attachment, metadata map[string]string) []map[string]interface{} {
	result := make([]map[string]interface{}, len(attachements))
	for i, attachement := range attachements {
		result[i] = map[string]interface{}{
			"id":            attachement.AttachmentID,
			"attached_at":   attachement.AttachedAt,
			"attached_mode": metadata["attached_mode"],
			"device_name":   attachement.Device,
			"server_id":     attachement.ServerID,
		}
	}
	return result
}

func sourceEvsVolumes(vols []volumes.Volume) ([]map[string]interface{}, []string, error) {
	result := make([]map[string]interface{}, len(vols))
	ids := make([]string, len(vols))

	for i, volume := range vols {
		vMap := map[string]interface{}{
			"id":                    volume.ID,
			"attachments":           sourceEvsAttachment(volume.Attachments, volume.Metadata),
			"availability_zone":     volume.AvailabilityZone,
			"description":           volume.Description,
			"volume_type":           volume.VolumeType,
			"enterprise_project_id": volume.EnterpriseProjectID,
			"name":                  volume.Name,
			"multiattach":           volume.Multiattach,
			"size":                  volume.Size,
			"status":                volume.Status,
			"created_at":            volume.CreatedAt,
			"updated_at":            volume.UpdatedAt,
			"wwn":                   volume.WWN,
			"metadata":              volume.Metadata,
			"tags":                  volume.Tags,
		}
		bootable, err := strconv.ParseBool(volume.Bootable)
		if err != nil {
			return nil, nil, fmtp.Errorf("The bootable of volume (%s) connot be converted from boolen to string.",
				volume.ID)
		} else {
			vMap["bootable"] = bootable
		}

		result[i] = vMap
		ids[i] = volume.ID
	}
	return result, ids, nil
}

func dataSourceEvsVolumesV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	client, err := cfg.BlockStorageV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloudStack EVS v2 client: %s", err)
	}

	pages, err := volumes.List(client, buildQueryOpts(d, cfg)).AllPages()
	if err != nil {
		return fmtp.DiagErrorf("An error occurred while fetching the pages of the EVS disks: %s", err)
	}
	vols, err := volumes.ExtractVolumes(pages)
	if err != nil {
		return fmtp.DiagErrorf("Error getting the EVS volume list form server: %s", err)
	}

	vMap, ids, err := sourceEvsVolumes(vols)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(hashcode.Strings(ids))
	if err = d.Set("volumes", vMap); err != nil {
		return fmtp.DiagErrorf("Error saving the detailed information of the EVS disks to state: %s", err)
	}
	return nil
}
