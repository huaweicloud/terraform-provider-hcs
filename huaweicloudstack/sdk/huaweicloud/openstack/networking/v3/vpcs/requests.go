package vpcs

import (
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

type UpdateOptsBuilder interface {
	ToVpcUpdateMap() (map[string]interface{}, error)
}

func AddSecondaryCIDR(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToVpcUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Put(addCidrURL(c, id), b, &r.Body, reqOpt)
	return
}

func RemoveSecondaryCIDR(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToVpcUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{200}}
	_, r.Err = c.Put(removeCidrURL(c, id), b, &r.Body, reqOpt)
	return
}

type UpdateOpts struct {
	ExtendCidrs []string `json:"extend_cidrs,omitempty"`
}

// ToVpcUpdateMap builds a create request body from CreateOpts.
func (opts UpdateOpts) ToVpcUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "vpc")
}

// GetVpcById is a method to obtain vpc informations from special region through vpc ID.
func Get(c *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(resourceURL(c, id), &r.Body, nil)
	return
}
