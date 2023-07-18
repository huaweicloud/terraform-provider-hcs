package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type LinksVoHasNextAndPre struct {
	Next *string `json:"next,omitempty"`

	Previous *string `json:"previous,omitempty"`

	Self *string `json:"self,omitempty"`
}

func (o LinksVoHasNextAndPre) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "LinksVoHasNextAndPre struct{}"
	}

	return strings.Join([]string{"LinksVoHasNextAndPre", string(data)}, " ")
}
