package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type TokenResponse struct {
	ExpiresAt *string `json:"expires_at,omitempty"`

	IssuedAt *string `json:"issued_at,omitempty"`

	Methods *[]string `json:"methods,omitempty"`

	Domain *NameIdVo `json:"domain,omitempty"`

	Project *UserVo `json:"project,omitempty"`

	Roles *[]NameIdVo `json:"roles,omitempty"`

	User *UserVo `json:"user,omitempty"`

	AssumedBy *UserVoWrapper `json:"assumed_by,omitempty"`

	Catalog *[]CatalogVo `json:"catalog,omitempty"`
}

func (o TokenResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TokenResponse struct{}"
	}

	return strings.Join([]string{"TokenResponse", string(data)}, " ")
}
