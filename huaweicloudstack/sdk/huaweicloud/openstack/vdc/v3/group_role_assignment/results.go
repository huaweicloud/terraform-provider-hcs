package group_role_assignment

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

type commonResult struct {
	golangsdk.Result
}

// User Group - Permission Management
type GroupRoleAssignmentListResult struct {
	commonResult
}

type CreateGroupRoleAssignmentResult struct {
	commonResult
}

type DeleteGroupRoleAssignmentResult struct {
	commonResult
}

type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Inherit     bool   `json:"inherit"`
	AssignType  string `json:"assign_type"`
	TargetID    string `json:"target_id"`
}

func (r GroupRoleAssignmentListResult) Extract() ([]Role, int, error) {
	var s struct {
		Roles []Role `json:"roles"`
		Total int    `json:"total"`
	}
	err := r.ExtractInto(&s)
	return s.Roles, s.Total, err
}
