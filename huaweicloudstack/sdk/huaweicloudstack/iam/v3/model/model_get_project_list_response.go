package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type GetProjectListResponse struct {
	Projects *[]ProjectDetailVo `json:"projects,omitempty"`

	Links          *LinksVoHasNextAndPre `json:"links,omitempty"`
	HttpStatusCode int                   `json:"-"`
}

func (o GetProjectListResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetProjectListResponse struct{}"
	}

	return strings.Join([]string{"GetProjectListResponse", string(data)}, " ")
}
