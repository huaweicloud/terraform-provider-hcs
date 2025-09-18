package direct_connects

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

// CreateOptsBuilder is an interface by which can be able to build the request body.
type CreateOptsBuilder interface {
	ToDcCreateMap() (map[string]interface{}, error)
}

// CreateOpts is a struct which represents the request body of create method.
type CreateOpts struct {
	Name         string   `json:"name,omitempty"`
	Type         string   `json:"type,omitempty"`
	SysTags      []string `json:"sys_tags,omitempty"`
	HostingId    string   `json:"hosting_id,omitempty"`
	PeerLocation string   `json:"peer_location,omitempty"`
	Description  string   `json:"description,omitempty"`
	Tenancy      string   `json:"tenancy,omitempty"`
}

func (opts CreateOpts) ToDcCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "direct_connect")
}

// Create is a method to create a new direct connect.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CommonResult) {
	b, err := opts.ToDcCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder is an interface by which can be able to build the request body.
type UpdateOptsBuilder interface {
	ToDcUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts is a struct which represents the request body of update method.
type UpdateOpts struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Tenancy     string `json:"tenancy,omitempty"`
}

func (opts UpdateOpts) ToDcUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "direct_connect")
}

// Update is a method to update an existing direct connect.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r CommonResult) {
	b, err := opts.ToDcUpdateMap()
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
	ToDcListQuery() (string, error)
}

// ListOpts allows extensions to add additional parameters to the API.
type ListOpts struct {
	ID       string `q:"id"`
	Name     string `q:"name"`
	Provider string `q:"provider"`
	Type     string `q:"type"`
}

// ToDcListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToDcListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List is a method by which can get the detailed information of all direct connects.
func List(c *golangsdk.ServiceClient, opts ListOptsBuilder) ([]DirectConnect, error) {
	url := rootURL(c)
	query, err := opts.ToDcListQuery()
	if err != nil {
		return nil, err
	}
	url += query
	var r CommonResult
	_, r.Err = c.Get(url, &r.Body, nil)

	return ExtractDirectConnects(r)
}

// Delete is a method to delete an existing direct connect.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}
