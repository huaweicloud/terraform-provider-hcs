package group_membership

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

type commonResult struct {
	golangsdk.Result
}

type GroupMembershipUserListResult struct {
	commonResult
}

type CreateGroupMembershipResult struct {
	commonResult
}

type UpdateGroupMembershipResult struct {
	commonResult
}

type DeleteGroupMembershipResult struct {
	commonResult
}

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

func (r GroupMembershipUserListResult) Extract() ([]User, int, error) {
	var s struct {
		Users []User `json:"users"`
		Total int    `json:"total"`
	}
	err := r.ExtractInto(&s)
	return s.Users, s.Total, err
}
