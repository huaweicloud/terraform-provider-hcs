package user

import golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF8", "X-Language": "en-us"},
}

var ApiAuthType map[string]string = map[string]string{
	"LOCAL_AUTH":   "0",
	"SAML_AUTH":    "1",
	"LDAP_AUTH":    "2",
	"MACHINE_USER": "4",
}

var ApiAccessMode map[string]string = map[string]string{
	"default":      "0",
	"console":      "1",
	"programmatic": "2",
}

type CreateOpts struct {
	Name        string `json:"name" required:"true"`
	Password    string `json:"password,omitempty" required:"false"`
	DisplayName string `json:"display_name,omitempty" required:"false"`
	AuthType    string `json:"auth_type" required:"false"`
	Enabled     bool   `json:"enabled" required:"false"`
	Description string `json:"description,omitempty" required:"false"`
	AccessMode  string `json:"access_mode" required:"false"`
}

func Create(client *golangsdk.ServiceClient, vdcId string, opts CreateOpts) (r CreatResult) {
	b, err := golangsdk.BuildRequestBody(opts, "user")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(listOrCreateURL(client, vdcId), b, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}

type PutOpts struct {
	DisplayName string `json:"display_name,omitempty" required:"false"`
	Enabled     bool   `json:"enabled" required:"false"`
	Description string `json:"description,omitempty" required:"false"`
	AccessMode  string `json:"access_mode" required:"false"`
}

func Update(client *golangsdk.ServiceClient, opts PutOpts, id string, isPwd bool) (r UpdateResult) {
	b, err := golangsdk.BuildRequestBody(opts, "user")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(detailOrPutURL(client, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{204},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}

type PutPwdOpts struct {
	Password string `json:"password,omitempty" required:"false"`
}

func UpPwd(client *golangsdk.ServiceClient, opts PutPwdOpts, id string) (r UpdateResult) {
	b, err := golangsdk.BuildRequestBody(opts, "user")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(pwdURL(client, id), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return
}

func Delete(client *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}

func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(detailOrPutURL(client, id), &r.Body, nil)
	return
}

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
	url := listOrCreateURL(c, vdcId)
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
