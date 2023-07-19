package model

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"

	"strings"
)

type GetDomainByUserTokenResponse struct {
	Domains *[]DomainResponse `json:"domains,omitempty"`

	Links          *LinksVoHasNextAndPre `json:"links,omitempty"`
	HttpStatusCode int                   `json:"-"`
}

func (o GetDomainByUserTokenResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetDomainByUserTokenResponse struct{}"
	}

	return strings.Join([]string{"GetDomainByUserTokenResponse", string(data)}, " ")
}
