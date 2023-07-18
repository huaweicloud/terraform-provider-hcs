package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type GetUserDetailRequest struct {
	XAuthToken string `json:"X-Auth-Token"`

	UserId string `json:"user_id"`
}

func (o GetUserDetailRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetUserDetailRequest struct{}"
	}

	return strings.Join([]string{"GetUserDetailRequest", string(data)}, " ")
}
