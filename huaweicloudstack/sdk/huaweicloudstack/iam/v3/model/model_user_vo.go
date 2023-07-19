package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type UserVo struct {
	Name *string `json:"name,omitempty"`

	Id *string `json:"id,omitempty"`

	Password *string `json:"password,omitempty"`

	Domain *NameIdVo `json:"domain,omitempty"`
}

func (o UserVo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UserVo struct{}"
	}

	return strings.Join([]string{"UserVo", string(data)}, " ")
}
