package v1peering

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

// ListOpts allows the filtering of collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the floating IP attributes you want to see returned.
type ListOpts struct {
	//ID is the unique identifier for the vpc_peering_connection.
	ID string `q:"id"`

	//Name is the human-readable name for the vpc_peering_connection. It does not have to be
	// unique.
	Name string `q:"name"`

	//Status indicates whether a vpc_peering_connection is currently operational.
	Status string `q:"status"`

	// VpcId indicates vpc_peering_connection available in specific vpc.
	VpcId string `q:"vpc_id"`
}

// List returns collection of vpc_peering_connection resources. It accepts
// a ListOpts struct, which allows you to filter the returned collection for greater efficiency.
func List(c *golangsdk.ServiceClient, opts ListOpts) (r GetResult) {
	q, err := golangsdk.BuildQueryString(&opts)
	if err != nil {
		return
	}
	u := rootURL(c) + q.String()
	_, r.Err = c.Get(u, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceListURL(c, id), &r.Body, nil)
	return
}

func Accept(c *golangsdk.ServiceClient, id string) (r AcceptResult) {
	_, r.Err = c.Put(acceptURL(c, id), nil, &r.Body, nil)
	return
}

// Reject is used by a tenant to reject a VPC peering connection request initiated by another tenant.
func Reject(c *golangsdk.ServiceClient, id string) (r RejectResult) {
	_, r.Err = c.Put(rejectURL(c, id), nil, &r.Body, nil)
	return
}

// CreateOptsBuilder is an interface by which can build the request body of vpc peering connection.
type CreateOptsBuilder interface {
	ToPeeringCreateMap() (map[string]interface{}, error)
}

// CreateOpts is a struct which is used to create vpc peering connection.
type CreateOpts struct {
	Name          string `json:"name"`
	LocalVpcId    string `json:"local_vpc_id" required:"true"`
	PeerVpcId     string `json:"peer_vpc_id" required:"true"`
	PeerRegion    string `json:"peer_region" required:"false"`
	PeerProjectId string `json:"peer_project_id" required:"false"`
}

// ToPeeringCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToPeeringCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create is a method by which can access to create the vpc peering connection.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToPeeringCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete is a method by which can be able to delete a vpc peering connection.
func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, id), nil)
	return
}

// UpdateOptsBuilder is an interface by which can be able to build the request body of vpc peering connection.
type UpdateOptsBuilder interface {
	ToVpcPeeringUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts is a struct which represents the request body of update method.
type UpdateOpts struct {
	Name string `json:"name,omitempty"`
}

// ToVpcPeeringUpdateMap builds an update request body from UpdateOpts.
func (opts UpdateOpts) ToVpcPeeringUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update is a method which can be able to update the name of vpc peering connection.
func Update(client *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToVpcPeeringUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(resourceURL(client, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
