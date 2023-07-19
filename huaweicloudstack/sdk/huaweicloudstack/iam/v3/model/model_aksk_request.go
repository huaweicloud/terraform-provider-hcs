package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type AkskRequest struct {
	Access *KeyVo `json:"access"`

	Secret *KeyVo `json:"secret"`
}

func (o AkskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AkskRequest struct{}"
	}

	return strings.Join([]string{"AkskRequest", string(data)}, " ")
}
