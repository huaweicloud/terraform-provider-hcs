package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type GetProjectDetailRequest struct {
	XAuthToken string `json:"X-Auth-Token"`

	ProjectId string `json:"project_id"`
}

func (o GetProjectDetailRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetProjectDetailRequest struct{}"
	}

	return strings.Join([]string{"GetProjectDetailRequest", string(data)}, " ")
}
