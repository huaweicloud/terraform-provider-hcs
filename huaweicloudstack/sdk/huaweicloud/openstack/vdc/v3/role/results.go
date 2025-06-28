package role

import golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

type VdcRoleResponse struct {
	Total       int            `json:"total"`
	SystemRoles []VdcRoleModel `json:"system_roles"`
	CustomRoles []VdcRoleModel `json:"custom_roles"`
}

type VdcRoleModel struct {
	ID            string     `json:"id"`
	DomainId      string     `json:"domain_id"`
	DomainName    string     `json:"domain_name"`
	Name          string     `json:"name"`
	DisplayName   string     `json:"display_name"`
	Flag          string     `json:"flag"`
	Catalog       string     `json:"catalog"`
	Type          string     `json:"type"`
	Description   string     `json:"description"`
	DescriptionCn string     `json:"description_cn"`
	CloudPlatform string     `json:"cloud_platform"`
	Policy        PolicyBase `json:"policy"`
	Tag           string     `json:"tag"`
	AppName       string     `json:"app_name"`
	DisplayType   string     `json:"display_type"`
}

type PolicyBase struct {
	Version   string           `json:"Version"`
	Depends   []VdcRoleDepends `json:"Depends"`
	Statement []StatementInfo  `json:"Statement"`
}

type VdcRoleDepends struct {
	Catalog     string `json:"Catalog"`
	DisplayName string `json:"Display_name"`
}

type StatementInfo struct {
	Effect   string      `json:"Effect"`
	Action   []string    `json:"Action"`
	Resource interface{} `json:"-"`
}

type ListResult struct {
	golangsdk.Result
}

func (listResult ListResult) Extract() ([]VdcRoleModel, int, error) {
	var result struct {
		SystemRoles []VdcRoleModel `json:"system_roles"`
		CustomRoles []VdcRoleModel `json:"custom_roles"`
		Total       int            `json:"total"`
	}
	err := listResult.Result.ExtractInto(&result)
	return append(result.SystemRoles, result.CustomRoles...), result.Total, err
}

type CreateResult struct {
	golangsdk.Result
}

func (createResult CreateResult) Extract() (*VdcRoleModel, error) {
	var result struct {
		Role *VdcRoleModel `json:"role"`
	}
	err := createResult.Result.ExtractInto(&result)
	return result.Role, err
}

type UpdateResult struct {
	golangsdk.Result
}

func (updateResult UpdateResult) Extract() (*VdcRoleModel, error) {
	var result struct {
		Role *VdcRoleModel `json:"role"`
	}
	err := updateResult.Result.ExtractInto(&result)
	return result.Role, err
}

type DeleteResult struct {
	golangsdk.ErrResult
}

type GetResult struct {
	golangsdk.Result
}

func (getResult GetResult) Extract() (*VdcRoleModel, error) {
	var result struct {
		Role *VdcRoleModel `json:"role"`
	}
	err := getResult.Result.ExtractInto(&result)
	return result.Role, err
}
