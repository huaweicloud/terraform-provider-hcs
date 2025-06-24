package group

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

type VdcUserGroupModel struct {
	Name        string `json:"name" required:"true"`
	Description string `json:"description,omitempty" required:"false"`
}

type CreateOpts struct {
	Group VdcUserGroupModel `json:"group" required:"true"`
}

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF8", "X-Language": "en-us"},
}

func Create(client *golangsdk.ServiceClient, vdcId string, opts CreateOpts) (r CreatResult) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(GroupURL(client, vdcId), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}

func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(UserGroupCommonURL(client, id), &r.Body, nil)
	return
}

type PutOpts struct {
	Name        string `json:"name" required:"false"`
	Description string `json:"description" required:"false"`
}

func Update(client *golangsdk.ServiceClient, opts PutOpts, id string) (r UpdateResult) {
	b, err := golangsdk.BuildRequestBody(opts, "group")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(UserGroupCommonURL(client, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{204},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}

func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(UserGroupCommonURL(client, id), nil)
	return
}

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
