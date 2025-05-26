package user

import golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const (
	Start   = 0
	Limit   = 100
	SortKey = "name"
	SortDir = "asc"
)

type ListOpts struct {
	Name    string `q:"name"`
	Start   int    `q:"start"`
	Limit   int    `q:"limit"`
	SortKey string `q:"sort_key"`
	SortDir string `q:"sort_dir"`
}

func (opts ListOpts) ToUserListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

type ListOptsBuilder interface {
	ToUserListQuery() (string, error)
}

func List(c *golangsdk.ServiceClient, vdcId string, opts ListOptsBuilder) (r ListResult) {
	url := rootURL(c, vdcId)
	if opts != nil {
		query, err := opts.ToUserListQuery()
		if err != nil {
			r.Err = err
		}
		url += query
	}

	_, r.Err = c.Get(url, &r.Body, nil)
	return
}
