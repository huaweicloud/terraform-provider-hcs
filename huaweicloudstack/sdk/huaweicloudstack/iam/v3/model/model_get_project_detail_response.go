package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type GetProjectDetailResponse struct {
	Project        *ProjectDetailVo `json:"project,omitempty"`
	HttpStatusCode int              `json:"-"`
}

func (o GetProjectDetailResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetProjectDetailResponse struct{}"
	}

	return strings.Join([]string{"GetProjectDetailResponse", string(data)}, " ")
}
