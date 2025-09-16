package vdc

import (
	"context"
	"fmt"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vdc/v3/group"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
)

func ResourceVdcUserGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVdcUserGroupCreate,
		ReadContext:   resourceVdcUserGroupRead,
		UpdateContext: resourceVdcUserGroupUpdate,
		DeleteContext: resourceVdcUserGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVdcUserGroupInstanceImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"vdc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceVdcUserGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	userGroupClient, err := hcsConfig.VdcClient(hcsConfig.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud VDC user group client %s", err)
	}

	vdcId := d.Get("vdc_id").(string)
	createOpts := group.CreateOpts{
		Group: group.VdcUserGroupModel{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		},
	}

	addUserGroup, err := group.Create(userGroupClient, vdcId, createOpts).ToExtract()

	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloudStack VDC user group: %s", err)
	}

	d.SetId(addUserGroup.ID)

	return resourceVdcUserGroupRead(ctx, d, meta)
}

func resourceVdcUserGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	userGroupClient, err := hcsConfig.VdcClient(hcsConfig.GetRegion(d))

	if err != nil {
		return fmtp.DiagErrorf("Unable to create HuaweiCloudStack VDC user group client : %s", err)
	}

	userGroupDetail, err := group.Get(userGroupClient, d.Id()).ToExtract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error retrieving HuaweiCloudStack VDC user group")
	}

	mErr := multierror.Append(nil,
		d.Set("name", userGroupDetail.Name),
		d.Set("description", userGroupDetail.Description),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error setting HuaweiCloudStack VDC user group fields: %w", err)
	}

	return nil
}

func resourceVdcUserGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChanges("vdc_id") {
		return diag.Errorf(`Unsupported attribute values for modification: "vdc_id".`)
	}

	hcsConfig := config.GetHcsConfig(meta)
	userGroupClient, err := hcsConfig.VdcClient(hcsConfig.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Unable to create HuaweiCloudStack VDC user group client : %s", err)
	}

	if d.HasChanges("description", "name") {
		updateOpts := group.PutOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		}

		_, err = group.Update(userGroupClient, updateOpts, d.Id()).ToExtract()
		if err != nil {
			return fmtp.DiagErrorf("Error updating HuaweiCloudStack VDC user group: %s", err)
		}
	}

	return resourceVdcUserGroupRead(ctx, d, meta)
}

func resourceVdcUserGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	userGroupClient, err := hcsConfig.VdcClient(hcsConfig.GetRegion(d))

	if err != nil {
		return fmtp.DiagErrorf("Unable to create HuaweiCloudStack VDC user group client : %s", err)
	}

	err = group.Delete(userGroupClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting VDC user group %s: %s", d.Id(), err)
	}

	d.SetId("")

	return nil
}

func resourceVdcUserGroupInstanceImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	hcsConfig := config.GetHcsConfig(meta)
	region := hcsConfig.GetRegion(d)
	userGroupClient, err := hcsConfig.VdcClient(region)
	if err != nil {
		return nil, fmt.Errorf("error creating VDC user group client: %s", err)
	}

	userGroupDetail, err := group.Get(userGroupClient, d.Id()).ToExtract()
	if err != nil {
		return nil, common.CheckDeleted(d, err, "VDC user group instance")
	}

	mErr := multierror.Append(nil,
		d.Set("vdc_id", userGroupDetail.VdcId),
		d.Set("name", userGroupDetail.Name),
		d.Set("description", userGroupDetail.Description),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmtp.Errorf("error setting HuaweiCloudStack VDC user group fields: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}
