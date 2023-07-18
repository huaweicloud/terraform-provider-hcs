package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type GetTokenResponse struct {
	Token          *TokenResponse `json:"token,omitempty"`
	HttpStatusCode int            `json:"-"`
}

func (o GetTokenResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetTokenResponse struct{}"
	}

	return strings.Join([]string{"GetTokenResponse", string(data)}, " ")
}
