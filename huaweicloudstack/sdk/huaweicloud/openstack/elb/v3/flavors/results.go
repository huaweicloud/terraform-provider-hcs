package flavors

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/pagination"
)

type Flavor struct {
	// Specifies the ID of the flavor.
	ID string `json:"id"`

	// Specifies the info of the flavor.
	Info FlavorInfo `json:"info"`

	// Specifies the name of the flavor.
	Name string `json:"name"`

	// Specifies whether shared.
	Shared bool `json:"shared"`

	// Specifies the type of the flavor.
	Type string `json:"type"`

	// Specifies whether sold out.
	SoldOut bool `json:"flavor_sold_out"`

	// Specifies whether bind.
	Status string `json:"status"`
}

type FlavorInfo struct {
	// Specifies the connection
	Connection *int `json:"connection"`

	// Specifies the cps.
	Cps *int `json:"cps"`

	// Specifies the qps
	Qps *int `json:"qps"`

	// Specifies the bandwidth
	Bandwidth *int `json:"bandwidth"`
}

// FlavorsPage is the page returned by a pager when traversing over a
// collection of flavor.
type FlavorsPage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks whether a FlavorsPage struct is empty.
func (r FlavorsPage) IsEmpty() (bool, error) {
	is, err := ExtractFlavors(r)
	return len(is) == 0, err
}

// ExtractFlavors accepts a Page struct, specifically a FlavorsPage struct,
// and extracts the elements into a slice of flavor structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractFlavors(r pagination.Page) ([]Flavor, error) {
	var s struct {
		Flavors []Flavor `json:"flavors"`
	}
	err := (r.(FlavorsPage)).ExtractInto(&s)
	return s.Flavors, err
}

type commonResult struct {
	golangsdk.Result
}

// Extract is a function that accepts a result and extracts a flavor.
func (r commonResult) Extract() (*Flavor, error) {
	var s struct {
		Flavor *Flavor `json:"flavor"`
	}
	err := r.ExtractInto(&s)
	return s.Flavor, err
}

// CreateResult represents the result of a Create operation. Call its Extract
// method to interpret the result as a flavor.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a Get operation. Call its Extract
// method to interpret the result as a flavor.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an Update operation. Call its Extract
// method to interpret the result as a flavor.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a Delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}
