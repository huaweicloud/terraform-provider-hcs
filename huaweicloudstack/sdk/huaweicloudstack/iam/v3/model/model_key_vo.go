package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type KeyVo struct {
	Key *string `json:"key,omitempty"`
}

func (o KeyVo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "KeyVo struct{}"
	}

	return strings.Join([]string{"KeyVo", string(data)}, " ")
}
