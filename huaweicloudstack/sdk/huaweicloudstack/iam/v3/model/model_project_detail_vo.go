package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type ProjectDetailVo struct {
	DomainId *string `json:"domain_id,omitempty"`

	IsDomain *bool `json:"is_domain,omitempty"`

	ParentId *string `json:"parent_id,omitempty"`

	Extra *interface{} `json:"extra,omitempty"`

	Description *string `json:"description,omitempty"`

	Name *string `json:"name,omitempty"`

	Id *string `json:"id,omitempty"`

	Enabled *bool `json:"enabled,omitempty"`

	Links *LinksVo `json:"links,omitempty"`
}

func (o ProjectDetailVo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProjectDetailVo struct{}"
	}

	return strings.Join([]string{"ProjectDetailVo", string(data)}, " ")
}
