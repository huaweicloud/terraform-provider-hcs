package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type RestrictVo struct {
	Impersonation *string `json:"impersonation,omitempty"`

	UserId *string `json:"user_id,omitempty"`

	UserName *string `json:"user_name,omitempty"`

	Roles *[]string `json:"roles,omitempty"`
}

func (o RestrictVo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RestrictVo struct{}"
	}

	return strings.Join([]string{"RestrictVo", string(data)}, " ")
}
