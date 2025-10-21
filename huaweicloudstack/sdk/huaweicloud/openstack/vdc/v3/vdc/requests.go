package vdc

import golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

type ListOpts struct {
	Name       string `q:"name"`
	UpperVdcId string `q:"upper_vdc_id"`
	QueryName  string `q:"query_name"`
}

func (opts ListOpts) ToVdcListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

type ListOptsBuilder interface {
	ToVdcListQuery() (string, error)
}

func List(c *golangsdk.ServiceClient, opts ListOptsBuilder) (r ListResult) {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToVdcListQuery()
		if err != nil {
			r.Err = err
		}
		url += query
	}

	_, r.Err = c.Get(url, &r.Body, nil)
	return
}
