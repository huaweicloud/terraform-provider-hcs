package domains

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to
// the List request
type ListOptsBuilder interface {
	ToDomainListQuery() (string, error)
}

// ListOpts provides options to filter the List results.
type ListOpts struct {
	// Enabled filters the response by enabled domains.
	Enabled *bool `q:"enabled"`

	// Name filters the response by domain name.
	Name string `q:"name"`
}

// ToDomainListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToDomainListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List enumerates the domains to which the current token has access.
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToDomainListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return DomainPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateOptsBuilder interface {
	ToDomainCreateMap() (map[string]interface{}, error)
}

// CreateOpts provides options used to create a domain.
type CreateOpts struct {
	// Name is the name of the new domain.
	Name string `json:"name" required:"true"`

	// Description is a description of the domain.
	Description string `json:"description,omitempty"`

	// Enabled sets the domain status to enabled or disabled.
	Enabled *bool `json:"enabled,omitempty"`
}

// ToDomainCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToDomainCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "domain")
}

// UpdateOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateOptsBuilder interface {
	ToDomainUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents parameters to update a domain.
type UpdateOpts struct {
	// Name is the name of the domain.
	Name string `json:"name,omitempty"`

	// Description is the description of the domain.
	Description string `json:"description,omitempty"`

	// Enabled sets the domain status to enabled or disabled.
	Enabled *bool `json:"enabled,omitempty"`
}

// ToUpdateCreateMap formats a UpdateOpts into an update request.
func (opts UpdateOpts) ToDomainUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "domain")
}
