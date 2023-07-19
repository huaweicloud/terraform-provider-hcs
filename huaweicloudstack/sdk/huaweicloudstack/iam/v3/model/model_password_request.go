package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type PasswordRequest struct {
	User *UserVo `json:"user"`
}

func (o PasswordRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PasswordRequest struct{}"
	}

	return strings.Join([]string{"PasswordRequest", string(data)}, " ")
}
