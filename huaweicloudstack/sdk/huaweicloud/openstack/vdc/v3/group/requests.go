package vdc

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

type ListReqParam struct {
	VdcID string `json:"vdc_id"`

	Name string `q:"name"`

	Start int `q:"start"`

	Limit int `q:"limit"`
}

func (opts ListReqParam) ToListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

func GetGroupList(c *golangsdk.ServiceClient, opts ListReqParam) (r ListResult) {
	url := GroupURL(c, opts.VdcID)
	query, err := opts.ToListQuery()
	if err != nil {
		r.Err = err
	}
	url += query

	_, r.Err = c.Get(url, &r.Body, nil)

	return
}
