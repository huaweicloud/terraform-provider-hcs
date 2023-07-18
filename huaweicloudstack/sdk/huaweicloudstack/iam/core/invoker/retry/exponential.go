package retry

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"
)

type Exponential struct {
}

func (e *Exponential) ComputeDelayBeforeNextRetry() int32 {
	return utils.Min32(MaxDelay, BaseDelay*(utils.Pow32(3, 2)))
}
