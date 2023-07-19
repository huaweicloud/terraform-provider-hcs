package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type GetProjectListRequest struct {
	XAuthToken string `json:"X-Auth-Token"`

	DomainId *string `json:"domain_id,omitempty"`

	Name *string `json:"name,omitempty"`

	ParentId *string `json:"parent_id,omitempty"`

	Enabled *string `json:"enabled,omitempty"`

	IsDomain *string `json:"is_domain,omitempty"`

	Page *string `json:"page,omitempty"`

	PerPage *string `json:"per_page,omitempty"`
}

func (o GetProjectListRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetProjectListRequest struct{}"
	}

	return strings.Join([]string{"GetProjectListRequest", string(data)}, " ")
}
