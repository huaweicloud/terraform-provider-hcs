package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type CatalogVo struct {
	Name *string `json:"name,omitempty"`

	Id *string `json:"id,omitempty"`

	Type *string `json:"type,omitempty"`

	Endpoints *[]EndpointVo `json:"endpoints,omitempty"`
}

func (o CatalogVo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CatalogVo struct{}"
	}

	return strings.Join([]string{"CatalogVo", string(data)}, " ")
}
