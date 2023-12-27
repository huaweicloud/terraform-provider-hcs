package volumetypes

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the List
// request.
type ListOptsBuilder interface {
	ToVolumeTypeListQuery() (string, error)
}

// ListOpts holds options for listing Volume Types. It is passed to the volumetypes.List
// function.
type ListOpts struct {
	// Comma-separated list of sort keys and optional sort directions in the
	// form of <key>[:<direction>].
	Sort string `q:"sort"`
	// Requests a page size of items.
	Limit int `q:"limit"`
	// Used in conjunction with limit to return a slice of items.
	Offset int `q:"offset"`

	ExtraSpecsOrigin map[string]interface{}

	// Arbitrary key-value pairs defined by the user.
	ExtraSpecs map[string]string `q:"extra_specs"`
}

// ToVolumeTypeListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToVolumeTypeListQuery() (string, error) {
	handleExtraSpecs(&opts)
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

func handleExtraSpecs(opts *ListOpts) {
	m := opts.ExtraSpecsOrigin
	if len(m) == 0 {
		return
	}
	nm := make(map[string]string, len(m))
	for k, v := range m {
		nm[k] = v.(string)
	}
	opts.ExtraSpecs = nm
}

// List returns Volume types.
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)

	if opts != nil {
		query, err := opts.ToVolumeTypeListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return VolumeTypePage{pagination.LinkedPageBase{PageResult: r}}
	})
}
