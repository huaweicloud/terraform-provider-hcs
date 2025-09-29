package direct_connects

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

type CommonResult struct {
	golangsdk.Result
}

type DeleteResult struct {
	golangsdk.ErrResult
}

// Extract is a function that accepts a result and extracts a DirectConnect resource.
func (r CommonResult) Extract() (*DirectConnect, error) {
	var s DirectConnect
	err := r.ExtractInto(&s)
	return &s, err
}

func (r CommonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "direct_connect")
}

// ExtractDirectConnects is a function that accepts a result and extracts a slice of DirectConnects structs.
func ExtractDirectConnects(r CommonResult) ([]DirectConnect, error) {
	var s []DirectConnect
	err := ExtractDirectConnectsInto(r, &s)
	return s, err
}

func ExtractDirectConnectsInto(r CommonResult, v interface{}) error {
	return r.Result.ExtractIntoSlicePtr(v, "direct_connects")
}

// DirectConnect is a struct that represents the detail of the direct connect.
type DirectConnect struct {
	// UUID for the direct connect.
	ID string `json:"id"`
	// Human-readable name for the direct connect. Might not be unique.
	Name string `json:"name"`

	// Indicates whether direct connect is currently operational. Possible values include
	// 'ACTIVE', 'DOWN', 'BUILD', 'ERROR', 'PENDING_DELETE', 'DELETED', 'APPLY', 'DENY', 'PENDING_PAY', 'PAID',
	// 'ORDERING', 'ACCEPT', or REJECTED.
	Status string `json:"status"`

	// User-defined description of the direct connect.
	Description string `json:"description"`

	// UUID of the hosting direct connect bound for the direct connect.
	HostingId string `json:"hosting_id"`

	// Provider of the direct connect, Possible values include: 'ce', 'vpc-gw'.
	Provider string `json:"provider"`

	// Type of the direct connect, Possible values include: 'hosted', 'hosting'.
	Type string `json:"type"`

	// User-defined peer location of the direct connect.
	PeerLocation string `json:"peer_location"`

	// Group of the direct connect.
	Group string `json:"group"`

	// Expiration time of the direct connect.
	Tenancy string `json:"tenancy"`
}
