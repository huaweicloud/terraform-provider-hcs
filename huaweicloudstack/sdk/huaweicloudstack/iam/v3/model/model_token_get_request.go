package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type TokenGetRequest struct {
	Identity *IdentityRequest `json:"identity"`

	Scope *ScopeRequest `json:"scope"`
}

func (o TokenGetRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TokenGetRequest struct{}"
	}

	return strings.Join([]string{"TokenGetRequest", string(data)}, " ")
}
