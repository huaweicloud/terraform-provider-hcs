package v1peeringroute

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

// ListOpts allows the filtering  of collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the floating IP attributes you want to see returned.
type ListOpts struct {
	PeeringId string `q:"peering_id"`
}

// List returns a collection of vpc_peering_connection  resources. It accepts
//a ListOpts struct, which allows you to filter  the returned collection for greater efficiency.
func List(c *golangsdk.ServiceClient, opts ListOpts, vpcId string) (r GetResult) {
	q, err := golangsdk.BuildQueryString(&opts)
	if err != nil {
		return
	}
	u := queryURL(c, vpcId) + q.String()
	_, r.Err = c.Get(u, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// OptsBuilder is an interface by which can build the request body of vpc peering connection.
type OptsBuilder interface {
	ToPeeringRouteList() ([]interface{}, error)
}

type CreateOpts struct {
	Route []Route `json:"routes"`
}

// ToPeeringCreateMap builds a create request body from CreateOpts.
func (opts CreateOpts) ToPeeringCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create is a method by which can access to create the vpc peering connection.
func Create(client *golangsdk.ServiceClient, opts CreateOpts, vpcId string) (r CreateResult) {
	b, err := opts.ToPeeringCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(addURL(client, vpcId), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete is a method by which can be able to delete a vpc peering connection.
func Delete(client *golangsdk.ServiceClient, opts CreateOpts, vpcId string) (r DeleteResult) {
	b, err := opts.ToPeeringCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.DeleteWithBody(removeURL(client, vpcId), b, nil)
	return
}
