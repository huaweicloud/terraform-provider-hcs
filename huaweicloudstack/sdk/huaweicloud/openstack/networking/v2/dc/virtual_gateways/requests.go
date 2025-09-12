package virtual_gateways

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

// CreateOptsBuilder is an interface by which can be able to build the request body.
type CreateOptsBuilder interface {
	ToVgwCreateMap() (map[string]interface{}, error)
}

// CreateOpts is a struct which represents the request body of create method.
type CreateOpts struct {
	Name        string         `json:"name,omitempty"`
	SysTags     []string       `json:"sys_tags,omitempty"`
	Description string         `json:"description,omitempty"`
	VpcGroup    []VpcGroupOpts `json:"vpc_group,omitempty"`
}

// VpcGroupOpts is a struct which represents the vpc group of create method body.
type VpcGroupOpts struct {
	VpcId          string `json:"vpc_id,omitempty"`
	LocalEpGroupId string `json:"local_ep_group_id,omitempty"`
}

func (opts CreateOpts) ToVgwCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "virtual_gateway")
}

// Create is a method to create a new virtual gateway.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CommonResult) {
	b, err := opts.ToVgwCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder is an interface by which can be able to build the request body.
type UpdateOptsBuilder interface {
	ToVgwUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts is a struct which represents the request body of update method.
type UpdateOpts struct {
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	VpcGroup    []VpcGroupOpts `json:"vpc_group,omitempty"`
}

func (opts UpdateOpts) ToVgwUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "virtual_gateway")
}

// Update is a method to update an existing virtual gateway.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r CommonResult) {
	b, err := opts.ToVgwUpdateMap()
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
	ToVgwListQuery() (string, error)
}

// ListOpts allows extensions to add additional parameters to the API.
type ListOpts struct {
	ID   string `q:"id"`
	Name string `q:"name"`
}

// ToVgwListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToVgwListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List is a method by which can get the detailed information of all virtual gateways.
func List(c *golangsdk.ServiceClient, opts ListOptsBuilder) ([]VirtualGateway, error) {
	url := rootURL(c)
	query, err := opts.ToVgwListQuery()
	if err != nil {
		return nil, err
	}
	url += query
	var r CommonResult
	_, r.Err = c.Get(url, &r.Body, nil)

	return ExtractVirtualGateways(r)
}

// Delete is a method to delete an existing virtual gateway.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}
