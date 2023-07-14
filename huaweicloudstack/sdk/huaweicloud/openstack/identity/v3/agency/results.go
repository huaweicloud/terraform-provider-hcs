package agency

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/identity/v3/roles"
)

type Agency struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	DomainID            string `json:"domain_id"`
	DelegatedDomainID   string `json:"trust_domain_id"`
	DelegatedDomainName string `json:"trust_domain_name"`
	Description         string `json:"description"`
	ExpireTime          string `json:"expire_time"`
	CreateTime          string `json:"create_time"`
	// validity period of the agency
	// In create and update response, it can be "FOREVER" or validity hours in int, for example, 480
	// In get response, it can be "FOREVER" or validity hours in string, for example, "480"
	Duration interface{} `json:"duration"`
}

type commonResult struct {
	golangsdk.Result
}

func (r commonResult) Extract() (*Agency, error) {
	var s struct {
		Agency *Agency `json:"agency"`
	}
	err := r.ExtractInto(&s)
	return s.Agency, err
}

type GetResult struct {
	commonResult
}

type CreateResult struct {
	commonResult
}

type UpdateResult struct {
	commonResult
}

type ErrResult struct {
	golangsdk.ErrResult
}

type ListRolesResult struct {
	golangsdk.Result
}

func (r ListRolesResult) ExtractRoles() ([]roles.Role, error) {
	var s struct {
		Roles []roles.Role `json:"roles"`
	}
	err := r.ExtractInto(&s)
	return s.Roles, err
}

type InheritedRole struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	// Links contains referencing links to the role.
	Links map[string]interface{} `json:"links"`
}

type ListInheritedRolesResult struct {
	golangsdk.Result
}

func (r ListInheritedRolesResult) ExtractRoles() ([]InheritedRole, error) {
	var s struct {
		Roles []InheritedRole `json:"roles"`
	}
	err := r.ExtractInto(&s)
	return s.Roles, err
}
