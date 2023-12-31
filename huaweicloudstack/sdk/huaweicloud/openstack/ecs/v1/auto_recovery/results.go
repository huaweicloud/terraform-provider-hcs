package auto_recovery

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

type AutoRecovery struct {
	SupportAutoRecovery string `json:"support_auto_recovery"`
}

type GetResult struct {
	golangsdk.Result
}

func (r GetResult) Extract() (*AutoRecovery, error) {
	s := &AutoRecovery{}
	return s, r.ExtractInto(s)
}
