package flavors

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToFlavorListMap() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	// Specifies the id.
	ID []string `q:"id"`
	// Specifies the name.
	Name []string `q:"name"`
	// Specifies whether shared.
	Shared *bool `q:"shared"`
	// Specifies the type.
	Type []string `q:"type"`
}

// ToFlavorListMap formats a ListOpts into a query string.
func (opts ListOpts) ToFlavorListMap() (string, error) {
	s, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return s.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// flavors.
func List(c *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		queryString, err := opts.ToFlavorListMap()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += queryString
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return FlavorsPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToFlavorCreateMap() (map[string]interface{}, error)
}

// Info specifies the details of a Flavor's QoS metrics.
type Info struct {
	FlavorType string `json:"flavor_type"` // Required. flavor (QoS) metric type: "cps", "connection", "bandwidth", ("qps" for l7 only).
	Value      int    `json:"value"`       // Required. Metric value; constrained by the cluster type and metric.
}

// CreateOpts is the common options struct used in this package's Create
// operation.
type CreateOpts struct {
	// Name of the flavor.
	Name string `json:"name,omitempty"`

	// ProjectID is the UUID of the project who owns the flavor.
	// Only administrative users can specify a project UUID other than their own.
	ProjectID string `json:"project_id,omitempty"`

	// Shared specifies whether the flavor is shared between projects.
	Shared bool `json:"shared,omitempty"`

	// Flavor (QoS) type. Only "l4" or "l7".
	Type string `json:"type,omitempty"`

	// TenantID is the UUID of the project who owns the flavor.
	// Only administrative users can specify a project UUID other than their own.
	TenantID string `json:"tenant_id,omitempty"`

	// List of Flavor (QoS) metric details.
	Info *[]Info `json:"info,omitempty"`
}

// ToFlavorCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToFlavorCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "flavor")
}

// Create accepts a CreateOpts struct and uses the values to create a new
// load balancer flavor.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToFlavorCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}

// Get retrieves a particular flavor based on its unique ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToFlavorUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts is the common options struct used in this package's Update
// operation.
type UpdateOpts struct {
	// Name of the flavor.
	Name *string `json:"name,omitempty"`

	// Updated QoS metrics list; replaces existing info.
	Info *[]Info `json:"info,omitempty"`
}

// ToFlavorUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToFlavorUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "flavor")
}

// Update allows flavor to be updated.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToFlavorUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will permanently delete a particular flavor based on its unique ID.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}
