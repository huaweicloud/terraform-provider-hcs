package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type ScopeRequest struct {
	Domain *NameIdVo `json:"domain,omitempty"`

	Project *NameIdVo `json:"project,omitempty"`
}

func (o ScopeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ScopeRequest struct{}"
	}

	return strings.Join([]string{"ScopeRequest", string(data)}, " ")
}
