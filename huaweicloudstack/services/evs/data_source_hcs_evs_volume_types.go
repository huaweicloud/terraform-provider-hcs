package evs

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/evs/volumetypes"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
)

func DataSourceEvsVolumeTypesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEvsVolumeTypesV2Read,

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"extra_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"volume_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_public": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extra_specs": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"qos_specs_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildQueryVolumeTypeOpts(d *schema.ResourceData) volumetypes.ListOpts {
	result := volumetypes.ListOpts{
		ExtraSpecsOrigin: d.Get("extra_specs").(map[string]interface{}),
	}
	return result
}

func sourceEvsVolumeTypes(vols []volumetypes.VolumeType) ([]map[string]interface{}, []string, error) {
	result := make([]map[string]interface{}, len(vols))
	ids := make([]string, len(vols))

	for i, volumeType := range vols {
		vMap := map[string]interface{}{
			"id":           volumeType.ID,
			"name":         volumeType.Name,
			"is_public":    volumeType.IsPublic,
			"description":  volumeType.Description,
			"extra_specs":  volumeType.ExtraSpecs,
			"qos_specs_id": volumeType.QosSpecID,
		}

		result[i] = vMap
		ids[i] = volumeType.ID
	}
	return result, ids, nil
}

func dataSourceEvsVolumeTypesV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	client, err := cfg.BlockStorageV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloudStack EVS v2 client: %s", err)
	}

	pages, err := volumetypes.List(client, buildQueryVolumeTypeOpts(d)).AllPages()
	if err != nil {
		return fmtp.DiagErrorf("An error occurred while fetching the pages of the EVS types: %s", err)
	}
	voltypes, err := volumetypes.ExtractVolumeTypes(pages)
	if err != nil {
		return fmtp.DiagErrorf("Error getting the EVS volume type list form server: %s", err)
	}

	filter := d.Get("availability_zone").(string)
	if filter != "" {
		filterVolumeTypes := filterVolumeTypeListByAz(voltypes, filter)
		voltypes = filterVolumeTypes
	}
	vTypeMap, ids, err := sourceEvsVolumeTypes(voltypes)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(hashcode.Strings(ids))
	if err = d.Set("volume_types", vTypeMap); err != nil {
		return fmtp.DiagErrorf("Error saving the detailed information of the EVS types to state: %s", err)
	}
	return nil
}

func filterVolumeTypeListByAz(volumeTypes []volumetypes.VolumeType, filterAz string) []volumetypes.VolumeType {
	result := make([]volumetypes.VolumeType, 0, len(volumeTypes))
	for _, volumeType := range volumeTypes {
		extraSpecs := volumeType.ExtraSpecs
		availabilityZone := extraSpecs["HW:availability_zone"]
		azList := strings.Split(availabilityZone, ",")
		for _, az := range azList {
			if az == filterAz {
				result = append(result, volumeType)
			}
		}
	}
	return result
}
