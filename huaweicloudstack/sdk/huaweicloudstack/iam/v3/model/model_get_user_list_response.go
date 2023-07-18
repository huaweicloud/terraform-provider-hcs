package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type GetUserListResponse struct {
	Users *[]UserDetailVo `json:"users,omitempty"`

	Links          *LinksVo `json:"links,omitempty"`
	HttpStatusCode int      `json:"-"`
}

func (o GetUserListResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetUserListResponse struct{}"
	}

	return strings.Join([]string{"GetUserListResponse", string(data)}, " ")
}
