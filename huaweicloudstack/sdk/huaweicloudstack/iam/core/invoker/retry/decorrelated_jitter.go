package retry

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"
)

type DecorRelatedJitter struct {
}

func (d *DecorRelatedJitter) ComputeDelayBeforeNextRetry() int32 {
	return utils.Min32(MaxDelay, utils.RandInt32(BaseDelay, 3*BaseDelay))
}
