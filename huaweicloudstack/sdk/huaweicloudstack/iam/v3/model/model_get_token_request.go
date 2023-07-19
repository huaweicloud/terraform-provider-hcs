package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type GetTokenRequest struct {
	XAuthToken string `json:"X-Auth-Token"`

	Nocatalog *string `json:"nocatalog,omitempty"`

	Body *TokenGetRequestBase `json:"body,omitempty"`
}

func (o GetTokenRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetTokenRequest struct{}"
	}

	return strings.Join([]string{"GetTokenRequest", string(data)}, " ")
}
