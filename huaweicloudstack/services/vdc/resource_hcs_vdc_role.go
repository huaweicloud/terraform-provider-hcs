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
	RoleSDK "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vdc/v3/role"
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
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
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

	domainId := configContext.Config.DomainID // 从全局配置中获取domain_id
	// 需要处理如下用户传入的参数
	// 用户传入domainId
	userInputDomainId := schemaResourceData.Get("domain_id").(string)
	if userInputDomainId != "" {
		domainId = userInputDomainId
	}

	policy := RoleSDK.PolicyBase{}
	policyDoc := schemaResourceData.Get("policy").(string)
	err = json.Unmarshal([]byte(policyDoc), &policy)
	if err != nil {
		return diag.Errorf("Error unmarshalling policy, please check the format of the policy document: %s", err)
	}
	createOpts := &RoleSDK.CreateOpts{
		DomainId: domainId, // 租户ID，租户侧用户调用时为必填参数。
		Role: RoleSDK.RequestBodyVdcRole{
			DisplayName: schemaResourceData.Get("name").(string),
			Type:        schemaResourceData.Get("type").(string),
			Description: schemaResourceData.Get("description").(string),
			Policy:      policy, // policy 用户传入policy内容json字符串
		},
	}
	role, err := RoleSDK.Create(vdcRoleClient, createOpts).Extract()
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

	role, err := RoleSDK.Get(vdcRoleClient, schemaResourceData.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(schemaResourceData, err, "IAM custom policy")
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
		schemaResourceData.Set("domain_id", role.DomainId),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("Error setting Vdc custom policy fields: %s", err)
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

	domainId := configContext.Config.DomainID // 从全局配置中获取domain_id
	// 需要处理如下用户传入的参数
	// 用户传入domainId
	userInputDomainId := schemaResourceData.Get("domain_id").(string)
	if userInputDomainId != "" {
		domainId = userInputDomainId
	}

	policy := RoleSDK.PolicyBase{}
	policyDoc := schemaResourceData.Get("policy").(string)
	err = json.Unmarshal([]byte(policyDoc), &policy)
	if err != nil {
		return diag.Errorf("Error unmarshalling policy, please check the format of the policy document: %s", err)
	}

	createOpts := &RoleSDK.CreateOpts{
		DomainId: domainId, // 租户ID，租户侧用户调用时为必填参数。
		Role: RoleSDK.RequestBodyVdcRole{
			DisplayName: schemaResourceData.Get("name").(string),
			Type:        schemaResourceData.Get("type").(string),
			Description: schemaResourceData.Get("description").(string),
			Policy:      policy, // policy 用户传入policy内容json字符串
		},
	}

	_, err = RoleSDK.Update(vdcRoleClient, schemaResourceData.Id(), createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating Vdc custom policy: %s", err)
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

	err = RoleSDK.Delete(vdcRoleClient, schemaResourceData.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("Error deleting Vdc custom policy: %s", err)
	}

	return nil
}
