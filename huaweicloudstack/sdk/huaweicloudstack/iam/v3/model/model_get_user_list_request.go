package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type GetUserListRequest struct {
	XAuthToken string `json:"X-Auth-Token"`

	DomainId *string `json:"domain_id,omitempty"`

	Enabled *string `json:"enabled,omitempty"`

	Name *string `json:"name,omitempty"`
}

func (o GetUserListRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetUserListRequest struct{}"
	}

	return strings.Join([]string{"GetUserListRequest", string(data)}, " ")
}
