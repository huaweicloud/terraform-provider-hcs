package dc_endpoint_groups

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

type CommonResult struct {
	golangsdk.Result
}

type DeleteResult struct {
	golangsdk.ErrResult
}

// Extract is a function that accepts a result and extracts a DcEndpointGroup resource.
func (r CommonResult) Extract() (*DcEndpointGroup, error) {
	var s DcEndpointGroup
	err := r.ExtractInto(&s)
	return &s, err
}

func (r CommonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "dc_endpoint_group")
}

// ExtractDcEndpointGroups is a function that accepts a result and extracts a slice of DcEndpointGroups structs.
func ExtractDcEndpointGroups(r CommonResult) ([]DcEndpointGroup, error) {
	var s []DcEndpointGroup
	err := ExtractDcEndpointGroupsInto(r, &s)
	return s, err
}

func ExtractDcEndpointGroupsInto(r CommonResult, v interface{}) error {
	return r.Result.ExtractIntoSlicePtr(v, "dc_endpoint_groups")
}

// DcEndpointGroup is a struct that represents the detail of the dc endpoint group.
type DcEndpointGroup struct {
	// UUID for the dc endpoint group.
	ID string `json:"id"`

	// Human-readable name for the dc endpoint group. Might not be unique.
	Name string `json:"name"`

	// User-defined description of the dc endpoint group.
	Description string `json:"description"`

	// Type of the dc endpoint group, only 'cidr' is allowed.
	Type string `json:"type"`

	// CIDR to be used.
	Endpoints []string `json:"endpoints"`
}
