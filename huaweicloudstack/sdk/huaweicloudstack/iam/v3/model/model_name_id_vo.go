package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type NameIdVo struct {
	Name *string `json:"name,omitempty"`

	Id *string `json:"id,omitempty"`
}

func (o NameIdVo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "NameIdVo struct{}"
	}

	return strings.Join([]string{"NameIdVo", string(data)}, " ")
}
