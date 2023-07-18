package retry

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"
)

type RandomJitter struct {
}

func (r *RandomJitter) ComputeDelayBeforeNextRetry() int32 {
	delay := utils.Min32(MaxDelay, BaseDelay*(utils.Pow32(3, 2)))
	return utils.RandInt32(0, delay)
}
