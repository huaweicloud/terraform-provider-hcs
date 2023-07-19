package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type EndpointVo struct {
	RegionId *string `json:"region_id,omitempty"`

	Id *string `json:"id,omitempty"`

	Region *string `json:"region,omitempty"`

	Interface *string `json:"interface,omitempty"`

	Url *string `json:"url,omitempty"`
}

func (o EndpointVo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EndpointVo struct{}"
	}

	return strings.Join([]string{"EndpointVo", string(data)}, " ")
}
