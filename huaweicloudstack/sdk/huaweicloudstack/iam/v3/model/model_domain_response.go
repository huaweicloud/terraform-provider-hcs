package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type DomainResponse struct {
	Name *string `json:"name,omitempty"`

	Description *string `json:"description,omitempty"`

	Id *string `json:"id,omitempty"`

	Enabled *bool `json:"enabled,omitempty"`

	Links *LinksVo `json:"links,omitempty"`
}

func (o DomainResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DomainResponse struct{}"
	}

	return strings.Join([]string{"DomainResponse", string(data)}, " ")
}
