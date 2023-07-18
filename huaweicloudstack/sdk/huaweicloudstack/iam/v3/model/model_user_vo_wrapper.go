package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type UserVoWrapper struct {
	User *UserVo `json:"user,omitempty"`
}

func (o UserVoWrapper) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UserVoWrapper struct{}"
	}

	return strings.Join([]string{"UserVoWrapper", string(data)}, " ")
}
