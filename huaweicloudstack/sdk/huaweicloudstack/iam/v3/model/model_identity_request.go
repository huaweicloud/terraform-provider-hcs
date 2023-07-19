package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type IdentityRequest struct {
	Methods []string `json:"methods"`

	Password *PasswordRequest `json:"password,omitempty"`

	HwAkSk *AkskRequest `json:"hw_ak_sk,omitempty"`

	HwAssumeRole *AssumeRequest `json:"hw_assume_role,omitempty"`
}

func (o IdentityRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "IdentityRequest struct{}"
	}

	return strings.Join([]string{"IdentityRequest", string(data)}, " ")
}
