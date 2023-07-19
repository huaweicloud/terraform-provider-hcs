package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type AssumeRequest struct {
	DomainName string `json:"domain_name"`

	XroleName string `json:"xrole_name"`

	Restrict *RestrictVo `json:"restrict,omitempty"`
}

func (o AssumeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AssumeRequest struct{}"
	}

	return strings.Join([]string{"AssumeRequest", string(data)}, " ")
}
