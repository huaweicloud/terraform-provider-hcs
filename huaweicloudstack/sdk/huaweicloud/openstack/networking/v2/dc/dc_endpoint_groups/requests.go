package dc_endpoint_groups

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

// CreateOptsBuilder is an interface by which can be able to build the request body.
type CreateOptsBuilder interface {
	ToDcEpGroupCreateMap() (map[string]interface{}, error)
}

// CreateOpts is a struct which represents the request body of create method.
type CreateOpts struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Type        string   `json:"type,omitempty"`
	Endpoints   []string `json:"endpoints,omitempty"`
}

func (opts CreateOpts) ToDcEpGroupCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "dc_endpoint_group")
}

// Create is a method to create a new dc endpoint group.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CommonResult) {
	b, err := opts.ToDcEpGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder is an interface by which can be able to build the request body.
type UpdateOptsBuilder interface {
	ToDcEpGroupUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts is a struct which represents the request body of update method.
type UpdateOpts struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (opts UpdateOpts) ToDcEpGroupUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "dc_endpoint_group")
}

// Update is a method to update an existing dc endpoint group.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r CommonResult) {
	b, err := opts.ToDcEpGroupUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, nil)
	return
}

// Get is a method to get the detailed information of a bandwidth.
func Get(c *golangsdk.ServiceClient, id string) (r CommonResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// ListOptsBuilder is an interface by which can be able to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToDcEpGroupListQuery() (string, error)
}

// ListOpts allows extensions to add additional parameters to the API.
type ListOpts struct {
	ID string `q:"id"`
}

// ToDcEpGroupListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToDcEpGroupListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List is a method by which can get the detailed information of all dc endpoint group.
func List(c *golangsdk.ServiceClient, opts ListOptsBuilder) (r CommonResult) {
	url := rootURL(c)
	query, err := opts.ToDcEpGroupListQuery()
	if err != nil {
		r.Err = err
		return
	}
	url += query
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// Delete is a method to delete an existing dc endpoint group.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}
