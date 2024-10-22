package v1peering

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"net/http"
)

type VpcInfo struct {
	DomainName string `json:"domain_name"`
	VpcId      string `json:"vpc_id"`
	TenantId   string `json:"tenant_id"`
}

// PeeringConnection represents a Neutron VPC peering connection.
// Manage and perform other operations on VPC peering connections,
// including querying VPC peering connections as well as
// creating, querying, deleting, and updating a VPC peering connection.
type PeeringConnection struct {
	// ID is the unique identifier for the vpc_peering_connection.
	ID string `json:"id"`

	// Name is the human-readable name for the vpc_peering_connection. It does not have to be
	// unique.
	Name string `json:"name"`

	// Status indicates whether a vpc_peering_connections is currently operational.
	RequesterVpcInfo VpcInfo `json:"requesterVpcInfo"`

	// Description is the supplementary information about the VPC peering connection.
	AccepterVpcInfo VpcInfo `json:"accepterVpcInfo"`

	// RequestVpcInfo indicates information about the local VPC
	Status string `json:"status"`
}

type commonResult struct {
	golangsdk.Result
}

// ExtractCreate is a function that accepts a result and extracts a Peering.
func (r commonResult) ExtractCreate() (PeeringConnection, error) {
	var s struct {
		Peering PeeringConnection `json:"vpc_peering_connection"`
	}
	err := r.ExtractInto(&s)
	return s.Peering, err
}

// ExtractList is a function that accepts a result and extracts a Peering.
func (r commonResult) ExtractList() ([]PeeringConnection, error) {
	var s struct {
		VpcPeeringConnections []PeeringConnection `json:"vpc_peering_connections"`
	}
	err1 := r.ExtractInto(&s)
	return s.VpcPeeringConnections, err1
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
	golangsdk.ErrResult
}

type AcceptResult struct {
	golangsdk.ErrResult
}

type Result struct {
	// Body is the payload of the HTTP response from the server. In most cases,
	// this will be the deserialized JSON structure.
	Body interface{}

	// Header contains the HTTP header structure from the original response.
	Header http.Header

	// Err is an error that occurred during the operation. It's deferred until
	// extraction to make it easier to chain the Extract call.
	Err error
}

type ErrResult struct {
	Result
}

func (r ErrResult) ExtractErr() error {
	return r.Err
}

// RejectResult represents the result of a get operation. Call its Extract
// method to interpret it as a Vpc Peering Connection.
type RejectResult struct {
	golangsdk.ErrResult
}
