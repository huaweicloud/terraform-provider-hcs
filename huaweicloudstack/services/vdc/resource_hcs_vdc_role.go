package vdc

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	roleSDK "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vdc/v3/role"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
)

// ResourceVdcRole
// @API VDC POST   /rest/vdc/v3.0/OS-ROLE/roles
// @API VDC GET    /rest/vdc/v3.0/OS-ROLE/roles/{role_id}
// @API VDC PUT    /rest/vdc/v3.0/OS-ROLE/roles/{role_id}
// @API VDC DELETE /rest/vdc/v3.0/OS-ROLE/roles/{role_id}
func ResourceVdcRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVdcRoleCreate,
		ReadContext:   resourceVdcRoleRead,
		UpdateContext: resourceVdcRoleUpdate,
		DeleteContext: resourceVdcRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					equal, _ := utils.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
		},
	}
}

func resourceVdcRoleCreate(ctx context.Context, schemaResourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configContext := config.GetHcsConfig(meta)
	region := configContext.GetRegion(schemaResourceData)
	vdcRoleClient, err := configContext.VdcClient(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating http client %s", err)
	}

	domainId := configContext.Config.DomainID

	policy := roleSDK.PolicyBase{}
	policyDoc := schemaResourceData.Get("policy").(string)
	err = json.Unmarshal([]byte(policyDoc), &policy)
	if err != nil {
		return diag.Errorf("Error unmarshalling policy, please check the format of the policy document: %s", err)
	}
	createOpts := &roleSDK.CreateOpts{
		DomainId: domainId,
		Role: roleSDK.RequestBodyVdcRole{
			DisplayName: schemaResourceData.Get("name").(string),
			Type:        schemaResourceData.Get("type").(string),
			Description: schemaResourceData.Get("description").(string),
			Policy:      policy,
		},
	}
	role, err := roleSDK.Create(vdcRoleClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating VDC custom role: %s", err)
	}

	schemaResourceData.SetId(role.ID)
	return resourceVdcRoleRead(ctx, schemaResourceData, meta)
}

func resourceVdcRoleRead(_ context.Context, schemaResourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configContext := config.GetHcsConfig(meta)
	region := configContext.GetRegion(schemaResourceData)
	vdcRoleClient, err := configContext.VdcClient(region)
	if err != nil {
		return diag.Errorf("Error creating http client: %s", err)
	}

	role, err := roleSDK.Get(vdcRoleClient, schemaResourceData.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(schemaResourceData, err, "VDC custom role")
	}

	policy, err := json.Marshal(role.Policy)
	if err != nil {
		return diag.Errorf("Error marshaling policy: %s", err)
	}

	mErr := multierror.Append(nil,
		schemaResourceData.Set("name", role.DisplayName),
		schemaResourceData.Set("description", role.Description),
		schemaResourceData.Set("type", role.Type),
		schemaResourceData.Set("policy", string(policy)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("Error setting VDC custom role fields: %s", err)
	}

	return nil
}
func resourceVdcRoleUpdate(ctx context.Context, schemaResourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configContext := config.GetHcsConfig(meta)
	region := configContext.GetRegion(schemaResourceData)
	vdcRoleClient, err := configContext.VdcClient(region)
	if err != nil {
		return diag.Errorf("Error creating http client: %s", err)
	}

	domainId := configContext.DomainID

	policy := roleSDK.PolicyBase{}
	policyDoc := schemaResourceData.Get("policy").(string)
	err = json.Unmarshal([]byte(policyDoc), &policy)
	if err != nil {
		return diag.Errorf("Error unmarshalling policy, please check the format of the policy document: %s", err)
	}

	if schemaResourceData.HasChanges("name", "description", "policy") {
		updateOpts := &roleSDK.UpdateOpts{
			DomainId: domainId,
			Role: roleSDK.RequestBodyVdcRole{
				DisplayName: schemaResourceData.Get("name").(string),
				Type:        schemaResourceData.Get("type").(string),
				Description: schemaResourceData.Get("description").(string),
				Policy:      policy,
			},
		}

		_, err = roleSDK.Update(vdcRoleClient, schemaResourceData.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating VDC custom role: %s", err)
		}
	}
	return resourceVdcRoleRead(ctx, schemaResourceData, meta)
}

func resourceVdcRoleDelete(_ context.Context, schemaResourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configContext := config.GetHcsConfig(meta)
	region := configContext.GetRegion(schemaResourceData)
	vdcRoleClient, err := configContext.VdcClient(region)
	if err != nil {
		return diag.Errorf("Error creating http client: %s", err)
	}

	err = roleSDK.Delete(vdcRoleClient, schemaResourceData.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("Error deleting VDC custom role: %s", err)
	}

	return nil
}
