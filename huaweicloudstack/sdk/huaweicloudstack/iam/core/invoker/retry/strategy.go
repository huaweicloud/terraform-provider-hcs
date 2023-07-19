package retry

const (
	BaseDelay = 5
	MaxDelay  = 60 * 1000
)

type Strategy interface {
	ComputeDelayBeforeNextRetry() int32
}
