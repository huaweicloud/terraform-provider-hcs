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
				Type:     schema.TypeString,
				Computed: true,
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
		return common.CheckDeletedDiag(schemaResourceData, err, "vdc project")
	}

	mErr := multierror.Append(nil,
		schemaResourceData.Set("name", project.Name),
		schemaResourceData.Set("display_name", project.DisplayName),
		schemaResourceData.Set("description", project.Description),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("Error setting Vdc project fields: %s", err)
	}

	return nil
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
			return diag.Errorf("Error updating Vdc project: %s", err)
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
		return diag.Errorf("Error deleting Vdc Vdc project: %s", err)
	}

	return nil
}
