package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type GetDomainByUserTokenRequest struct {
	XAuthToken string `json:"X-Auth-Token"`
}

func (o GetDomainByUserTokenRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetDomainByUserTokenRequest struct{}"
	}

	return strings.Join([]string{"GetDomainByUserTokenRequest", string(data)}, " ")
}
