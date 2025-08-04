package vdc

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	sdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vdc/v3/project"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
	"strings"
)

// ResourceVdcProject
// @API VDC POST /rest/vdc/v3.1/projects
// @API VDC GET /rest/vdc/v3.1/projects/{project_id}
// @API VDC PUT /rest/vdc/v3.0/projects/{project_id}
// @API VDC DELETE /rest/vdc/v3.0/projects/{project_id}
func ResourceVdcProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVdcProjectCreate,
		ReadContext:   resourceVdcProjectRead,
		UpdateContext: resourceVdcProjectUpdate,
		DeleteContext: resourceVdcProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: utils.ValidateVdcProjectInputName,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vdc_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     regionsSchema(),
			},
			"is_support_hws_service": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_shared": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func regionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_name": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"region_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_infras": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     cloudInfrasSchema(),
			},
		},
	}
}

func cloudInfrasSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cloud_infra_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_infra_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_infra_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_infra_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"azs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     azsSchema(),
			},
			"quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     quotasSchema(),
			},
		},
	}
}

func azsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"az_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"az_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"az_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func quotasSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cloud_infra_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_infra_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_infra_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_infra_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"azs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     azsSchema(),
			},
		},
	}
}

func resourceVdcProjectCreate(ctx context.Context, schemaResourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configContext := config.GetHcsConfig(meta)
	region := configContext.GetRegion(schemaResourceData)
	vdcClient, err := configContext.VdcClient(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating http client %s", err)
	}

	name := schemaResourceData.Get("name").(string)
	temps := strings.Split(name, "_")

	if len(temps) < 2 {
		return fmtp.DiagErrorf("Error the format of the name input parameter is incorrect. It should start with \"${region_id}_\" .  %s", name)
	}

	regionId := temps[0]
	createOpts := &sdk.CreateOpts{
		Project: sdk.CreateProjectV31{
			Name:                schemaResourceData.Get("name").(string),
			TenantId:            schemaResourceData.Get("vdc_id").(string),
			Description:         schemaResourceData.Get("description").(string),
			Regions:             []sdk.Region{{RegionId: regionId}},
			IsSupportHwsService: true,
		},
	}
	project, err := sdk.Create(vdcClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating VDC project: %s", err)
	}

	schemaResourceData.SetId(project.Id)
	return resourceVdcProjectRead(ctx, schemaResourceData, meta)
}

func resourceVdcProjectRead(_ context.Context, schemaResourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configContext := config.GetHcsConfig(meta)
	region := configContext.GetRegion(schemaResourceData)
	vdcClient, err := configContext.VdcClient(region)
	if err != nil {
		return diag.Errorf("Error creating http client: %s", err)
	}

	getOpts := &sdk.GetOpts{
		ServiceId: "",
	}

	project, err := sdk.Get(vdcClient, schemaResourceData.Id(), getOpts).Extract()
	if err != nil {
		return common.CheckDeletedDiag(schemaResourceData, err, "VDC project")
	}

	elements := getRegionElements(*project)

	mErr := multierror.Append(nil,
		schemaResourceData.Set("name", project.Name),
		schemaResourceData.Set("display_name", project.DisplayName),
		schemaResourceData.Set("description", project.Description),
		schemaResourceData.Set("is_support_hws_service", project.IsSupportHwsService),
		schemaResourceData.Set("is_shared", project.IsShared),
		schemaResourceData.Set("regions", elements),
		schemaResourceData.Set("vdc_id", project.TenantId),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("Error setting VDC project fields: %s", err)
	}

	return nil
}

func getRegionElements(project sdk.QueryProjectDetailV31) []map[string]interface{} {
	elements := make([]map[string]interface{}, 0, len(project.Regions))
	for _, group := range project.Regions {
		groupMap := map[string]interface{}{
			"region_id":     group.RegionId,
			"region_name":   getRegionName(group.RegionName),
			"region_type":   group.RegionType,
			"region_status": group.RegionStatus,
			"cloud_infras":  getCloudInfras(group.CloudInfras),
		}

		elements = append(elements, groupMap)
	}
	return elements
}

func getRegionName(regionName sdk.RegionName) map[string]string {
	element := make(map[string]string)
	element["en_us"] = regionName.EnUs
	element["zh_cn"] = regionName.ZhCn
	return element
}

func getCloudInfras(cloudInfras []sdk.ProjectCloudInfraInfo) []map[string]interface{} {
	elements := make([]map[string]interface{}, 0, len(cloudInfras))
	for _, group := range cloudInfras {
		groupMap := map[string]interface{}{
			"cloud_infra_id":     group.CloudInfraId,
			"cloud_infra_name":   group.CloudInfraName,
			"cloud_infra_type":   group.CloudInfraType,
			"cloud_infra_status": group.CloudInfraStatus,
			"azs":                getAzs(group.Azs),
			"quotas":             getQuotas(group.Quotas),
		}

		elements = append(elements, groupMap)
	}
	return elements
}

func getAzs(azs []sdk.ProjectAzInfo) []map[string]interface{} {
	elements := make([]map[string]interface{}, 0, len(azs))
	for _, group := range azs {
		groupMap := map[string]interface{}{
			"az_id":     group.AzId,
			"az_name":   group.AzName,
			"az_status": group.AzStatus,
		}

		elements = append(elements, groupMap)
	}
	return elements
}

func getQuotas(quotas []sdk.QuotaInfo) []map[string]interface{} {
	elements := make([]map[string]interface{}, 0, len(quotas))
	for _, group := range quotas {
		groupMap := map[string]interface{}{
			"service_id": group.ServiceId,
			"action":     group.Action,
			"resources":  getDictionaryList(group.Resources),
		}

		elements = append(elements, groupMap)
	}
	return elements
}

func getDictionaryList(dictList []sdk.DictInfo) []map[string]interface{} {
	elements := make([]map[string]interface{}, 0, len(dictList))
	for _, group := range dictList {
		groupMap := map[string]interface{}{
			"resource":    group.Resource,
			"local_limit": group.LocalLimit,
			"other_limit": group.OtherLimit,
			"local_used":  group.LocalUsed,
			"other_used":  group.OtherUsed,
		}

		elements = append(elements, groupMap)
	}
	return elements
}

func resourceVdcProjectUpdate(ctx context.Context, schemaResourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configContext := config.GetHcsConfig(meta)
	region := configContext.GetRegion(schemaResourceData)
	vdcClient, err := configContext.VdcClient(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating http client %s", err)
	}

	name := schemaResourceData.Get("name").(string)
	temps := strings.Split(name, "_")

	if len(temps) < 2 {
		return fmtp.DiagErrorf("Error the format of the name input parameter is incorrect.  %s", name)
	}

	if schemaResourceData.HasChanges("vdc_id") {
		return fmtp.DiagErrorf("Unsupported attribute values for modification: \"vdc_id\".")
	}

	if schemaResourceData.HasChanges("name", "description", "display_name") {
		updateOpts := &sdk.UpdateOpts{
			Project: sdk.UpdateProject{
				Name:        schemaResourceData.Get("name").(string),
				DisplayName: schemaResourceData.Get("display_name").(string),
				Description: schemaResourceData.Get("description").(string),
			},
		}

		err = sdk.Update(vdcClient, schemaResourceData.Id(), updateOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("Error updating VDC project: %s", err)
		}
	}

	return resourceVdcProjectRead(ctx, schemaResourceData, meta)
}

func resourceVdcProjectDelete(_ context.Context, schemaResourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configContext := config.GetHcsConfig(meta)
	region := configContext.GetRegion(schemaResourceData)
	vdcClient, err := configContext.VdcClient(region)
	if err != nil {
		return diag.Errorf("Error creating http client: %s", err)
	}
	deleteOpts := &sdk.DeleteOpts{
		DeleteServiceTypes: "",
	}
	err = sdk.Delete(vdcClient, schemaResourceData.Id(), deleteOpts).ExtractErr()
	if err != nil {
		return diag.Errorf("Error deleting VDC project: %s", err)
	}

	return nil
}
