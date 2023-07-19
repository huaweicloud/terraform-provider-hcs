package retry

type None struct {
}

func (n *None) ComputeDelayBeforeNextRetry() int32 {
	return 0
}
