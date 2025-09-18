package virtual_gateways

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

type CommonResult struct {
	golangsdk.Result
}

type DeleteResult struct {
	golangsdk.ErrResult
}

// Extract is a function that accepts a result and extracts a VirtualGateway resource.
func (r CommonResult) Extract() (*VirtualGateway, error) {
	var s VirtualGateway
	err := r.ExtractInto(&s)
	return &s, err
}

func (r CommonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "virtual_gateway")
}

// ExtractVirtualGateways is a function that accepts a result and extracts a slice of VirtualGateways structs.
func ExtractVirtualGateways(r CommonResult) ([]VirtualGateway, error) {
	var s []VirtualGateway
	err := ExtractVirtualGatewaysInto(r, &s)
	return s, err
}

func ExtractVirtualGatewaysInto(r CommonResult, v interface{}) error {
	return r.Result.ExtractIntoSlicePtr(v, "virtual_gateways")
}

// VirtualGateway is a struct that represents the detail of the virtual gateway.
type VirtualGateway struct {
	// UUID for the virtual gateway.
	ID string `json:"id"`
	// Human-readable name for the virtual gateway. Might not be unique.
	Name string `json:"name"`

	// Indicates whether virtual gateway is currently operational. Possible values include
	// 'ACTIVE', 'DOWN', 'BUILD', 'ERROR', 'PENDING_CREATE', 'PENDING_UPDATE', 'PENDING_DELETE'.
	Status string `json:"status"`

	// User-defined description of the virtual gateway.
	Description string `json:"description"`

	// Vpc info bound with the virtual gateway.
	VpcGroup []VpcGroup `json:"vpc_group"`
}

type VpcGroup struct {
	// UUID of the vpc.
	VpcId string `json:"vpc_id"`

	// UUID of the dc endpoint group to be bound.
	LocalEpGroupId string `json:"local_ep_group_id"`
}
