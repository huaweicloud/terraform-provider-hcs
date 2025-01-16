package obs

const (
	// SubResourceCustomDomain subResource value: customdomain
	SubResourceCustomDomain SubResourceType = "customdomain"
)

// RedundancyType defines type of redundancyType
type BucketRedundancyType string

const (
	BucketRedundancyClassic BucketRedundancyType = "CLASSIC"
	BucketRedundancyFusion  BucketRedundancyType = "FUSION"
)
