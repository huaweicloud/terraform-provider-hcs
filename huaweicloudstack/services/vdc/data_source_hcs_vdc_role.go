package vdc

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	RoleSDK "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vdc/v3/role"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
)

// DataSourceVdcRole @API VDC GET /rest/vdc/v3.0/OS-ROLE/roles/third-party/roles
func DataSourceVdcRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVdcRoleRead,

		Schema: map[string]*schema.Schema{
			"role_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"name", "display_name"},
			},
			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"name", "display_name"},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"catalog": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

const (
	RoleTypeCustom string = "custom"
	RoleTypeSystem string = "system"
)

func dataSourceVdcRoleRead(_ context.Context, schemaResourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configContext := config.GetHcsConfig(meta)
	region := configContext.GetRegion(schemaResourceData)
	vdcRoleClient, err := configContext.VdcClient(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating http client %s", err)
	}

	domainId := configContext.Config.DomainID // 从全局配置中获取domain_id

	var isSystem string
	//role_type 用户传入查询类型参数，根据类型过滤, system, custom
	userInputRoleType := schemaResourceData.Get("role_type").(string)
	existRoleType := userInputRoleType != ""

	if existRoleType {
		if userInputRoleType == RoleTypeSystem {
			isSystem = "true"
		} else if userInputRoleType == RoleTypeCustom {
			isSystem = "false"
		} else {
			return fmtp.DiagErrorf("Error retrieving optional parameter role_type %s, just support value %s or %s.", userInputRoleType, RoleTypeSystem, RoleTypeCustom)
		}
	}
	// display_name 用户传入显示名参数 根据显示名精确查找
	userInputDisplayName := schemaResourceData.Get("display_name").(string)
	existDisplayName := userInputDisplayName != ""

	// name 用户传入名称参数， 根据name精确查找
	userInputName := schemaResourceData.Get("name").(string)
	existName := userInputName != ""

	listOpts := RoleSDK.ListOpts{
		DomainId:    domainId, // 租户ID，租户侧用户调用时为必填参数，管理侧用户调用时为选填参数。
		IsSystem:    isSystem, // is_system=true：系统角色+系统策略， is_system=false：自定义策略，不传表示查询所有。
		FineGrained: false,    // 是否支持细粒度策略，不包含云管角色。
		Start:       0,        // 分页查询的起始位置，取值在0-2147483647之间，默认从0开始。
		Limit:       100,      // 限制每页显示的条目数量，取值在0-100之间。
	}

	var roles []RoleSDK.VdcRoleModel

	roles, err = findVdcRoleList(vdcRoleClient, listOpts, roles, func(tempAllRoles []RoleSDK.VdcRoleModel) ([]RoleSDK.VdcRoleModel, bool) {
		var exist bool
		var findResultRoles = make([]RoleSDK.VdcRoleModel, len(tempAllRoles))
		copy(findResultRoles, tempAllRoles)
		if existName {
			findResultRoles, exist = findRoleByName(findResultRoles, userInputName)
		}
		if existDisplayName {
			findResultRoles, exist = findRoleByDisplayName(findResultRoles, userInputDisplayName)
		}
		return findResultRoles, exist
	})

	if err != nil {
		return fmtp.DiagErrorf("Error retrieving VDC roles %s", err)
	}

	if len(roles) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(roles) > 1 {
		return diag.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria.")
	}

	role := roles[0]

	schemaResourceData.SetId(role.ID)
	return dataSourceVdcRoleAttributes(schemaResourceData, &role)
}

func dataSourceVdcRoleAttributes(schemaResourceData *schema.ResourceData, role *RoleSDK.VdcRoleModel) diag.Diagnostics {
	policy, err := json.Marshal(role.Policy)
	if err != nil {
		return diag.Errorf("Error marshaling the policy of VDC role: %s", err)
	}

	mErr := multierror.Append(nil,
		schemaResourceData.Set("domain_id", role.DomainId),
		schemaResourceData.Set("name", role.Name),
		schemaResourceData.Set("description", role.Description),
		schemaResourceData.Set("display_name", role.DisplayName),
		schemaResourceData.Set("type", role.Type),
		schemaResourceData.Set("policy", string(policy)),
		schemaResourceData.Set("catalog", role.Catalog),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("Error setting VDC role fields: %s", err)
	}
	return nil
}

// findVdcRoleList 用于调用接口获取列表数据
func findVdcRoleList(vdcRoleClient *golangsdk.ServiceClient, listOpts RoleSDK.ListOpts, tempAllRoles []RoleSDK.VdcRoleModel, callback func([]RoleSDK.VdcRoleModel) ([]RoleSDK.VdcRoleModel, bool)) ([]RoleSDK.VdcRoleModel, error) {
	var allRoles = make([]RoleSDK.VdcRoleModel, len(tempAllRoles))
	copy(allRoles, tempAllRoles)
	for { // 开启轮询

		// 查询一次列表数据
		vdcRoleResponse, total, err := RoleSDK.List(vdcRoleClient, listOpts).Extract()

		if err != nil {
			fmtp.DiagErrorf("Unable to query vdc roles: %s", err)
			return []RoleSDK.VdcRoleModel{}, err
		}
		for _, item := range vdcRoleResponse {
			allRoles = append(allRoles, item)
		}

		// 基于当前查询的查找目标元素，如果查到，就返回目标元素
		targetRoles, ok := callback(vdcRoleResponse)
		if ok {
			// 查到了就返回最终结果
			return targetRoles, nil
		}
		// 未查找到目标元素， 当前累计查询的数量还小于总数
		if total > len(allRoles) { // 更新查询条件，等待下一次轮询
			listOpts = RoleSDK.ListOpts{
				DomainId:    listOpts.DomainId,    // 租户ID，租户侧用户调用时为必填参数，管理侧用户调用时为选填参数。
				IsSystem:    listOpts.IsSystem,    // is_system=true：系统角色+系统策略， is_system=false：自定义策略，不传表示查询所有。
				FineGrained: listOpts.FineGrained, // 是否支持细粒度策略，不包含云管角色。
				Start:       len(allRoles),        // 分页查询的起始位置，取值在0-2147483647之间，默认从0开始。
				Limit:       listOpts.Limit,       // 限制每页显示的条目数量，取值在0-100之间。
			}
		} else {
			break
		}
	}
	// 没有查询到时，轮询结束，返回一个空结果
	return []RoleSDK.VdcRoleModel{}, nil
}

func findRoleByName(allVdcRoleResponseList []RoleSDK.VdcRoleModel, findValue string) ([]RoleSDK.VdcRoleModel, bool) {
	for _, roleModel := range allVdcRoleResponseList {
		value := roleModel.Name

		if value == findValue {
			return []RoleSDK.VdcRoleModel{roleModel}, true
		}
	}
	return []RoleSDK.VdcRoleModel{}, false
}

func findRoleByDisplayName(allVdcRoleResponseList []RoleSDK.VdcRoleModel, findValue string) ([]RoleSDK.VdcRoleModel, bool) {
	for _, roleModel := range allVdcRoleResponseList {
		value := roleModel.DisplayName

		if value == findValue {
			return []RoleSDK.VdcRoleModel{roleModel}, true
		}
	}
	return []RoleSDK.VdcRoleModel{}, false
}
