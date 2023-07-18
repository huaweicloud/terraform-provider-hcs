package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type GetUserDetailResponse struct {
	User           *UserDetailVo `json:"user,omitempty"`
	HttpStatusCode int           `json:"-"`
}

func (o GetUserDetailResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetUserDetailResponse struct{}"
	}

	return strings.Join([]string{"GetUserDetailResponse", string(data)}, " ")
}
