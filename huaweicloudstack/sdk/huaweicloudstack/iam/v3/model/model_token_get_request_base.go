package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type TokenGetRequestBase struct {
	Auth *TokenGetRequest `json:"auth"`
}

func (o TokenGetRequestBase) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TokenGetRequestBase struct{}"
	}

	return strings.Join([]string{"TokenGetRequestBase", string(data)}, " ")
}
