package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type LinksVo struct {
	Self *string `json:"self,omitempty"`
}

func (o LinksVo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "LinksVo struct{}"
	}

	return strings.Join([]string{"LinksVo", string(data)}, " ")
}
