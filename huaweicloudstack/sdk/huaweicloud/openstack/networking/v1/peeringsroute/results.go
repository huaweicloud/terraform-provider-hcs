package v1peeringroute

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

type Route struct {
	Destination string `json:"destination"`
	NextHop     string `json:"nexthop"`
}

type Response struct {
	routes []Route `json:"routes"`
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a Peering.
func (r commonResult) Extract() ([]Route, error) {
	s := struct {
		Route []Route `json:"routes"`
	}{}
	err := r.ExtractInto(&s)
	return s.Route, err
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Vpc Peering Connection.
type GetResult struct {
	commonResult
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a vpc peering connection.
type CreateResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a vpc peering connection.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	commonResult
}
