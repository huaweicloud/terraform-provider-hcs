package retry

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"
)

// EqualJitter 等抖动指数退避 delay = Exponential/2 + random(0, Exponential/2)
type EqualJitter struct {
}

func (e *EqualJitter) ComputeDelayBeforeNextRetry() int32 {
	delay := utils.Min32(MaxDelay, BaseDelay*(utils.Pow32(3, 2)))
	return delay/2 + utils.RandInt32(0, delay/2)
}
