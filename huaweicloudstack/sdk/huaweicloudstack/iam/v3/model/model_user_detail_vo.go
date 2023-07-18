package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type UserDetailVo struct {
	DomainId *string `json:"domain_id,omitempty"`

	Description *string `json:"description,omitempty"`

	Name *string `json:"name,omitempty"`

	PasswordExpiresAt *string `json:"password_expires_at,omitempty"`

	Id *string `json:"id,omitempty"`

	Enabled *bool `json:"enabled,omitempty"`

	Links *LinksVo `json:"links,omitempty"`
}

func (o UserDetailVo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UserDetailVo struct{}"
	}

	return strings.Join([]string{"UserDetailVo", string(data)}, " ")
}
