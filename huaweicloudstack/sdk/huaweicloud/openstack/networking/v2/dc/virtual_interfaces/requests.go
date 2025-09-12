package virtual_interfaces

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

// CreateOptsBuilder is an interface by which can be able to build the request body.
type CreateOptsBuilder interface {
	ToVifCreateMap() (map[string]interface{}, error)
}

// CreateOpts is a struct which represents the request body of create method.
type CreateOpts struct {
	Name            string         `json:"name,omitempty"`
	SysTags         []string       `json:"sys_tags,omitempty"`
	DirectConnectId string         `json:"direct_connect_id,omitempty"`
	VgwId           string         `json:"vgw_id,omitempty"`
	RemoteEpGroupId string         `json:"remote_ep_group_id,omitempty"`
	Description     string         `json:"description,omitempty"`
	LinkInfos       []LinkInfoOpts `json:"link_infos,omitempty"`
}

// LinkInfoOpts is a struct which represents the request query of get method.
type LinkInfoOpts struct {
	InterfaceGroupId  string `json:"interface_group_id,omitempty"`
	HostingId         string `json:"hosting_id,omitempty"`
	LocalGatewayV4Ip  string `json:"local_gateway_v4_ip,omitempty"`
	LocalGatewayV6Ip  string `json:"local_gateway_v6_ip,omitempty"`
	RemoteGatewayV4Ip string `json:"remote_gateway_v4_ip,omitempty"`
	RemoteGatewayV6Ip string `json:"remote_gateway_v6_ip,omitempty"`
	Vlan              int    `json:"vlan,omitempty"`
	BgpAsn            int    `json:"bgp_asn,omitempty"`
	BgpAsnDot         string `json:"bgp_asn_dot,omitempty"`
}

func (opts CreateOpts) ToVifCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "virtual_interface")
}

// Create is a method to create a new virtual interface.
func Create(c *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CommonResult) {
	b, err := opts.ToVifCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(rootURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder is an interface by which can be able to build the request body.
type UpdateOptsBuilder interface {
	ToVifUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts is a struct which represents the request body of update method.
type UpdateOpts struct {
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	RemoteEpGroupId string `json:"remote_ep_group_id,omitempty"`
}

func (opts UpdateOpts) ToVifUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "virtual_interface")
}

// Update is a method to update an existing virtual interface.
func Update(c *golangsdk.ServiceClient, id string, opts UpdateOptsBuilder) (r CommonResult) {
	b, err := opts.ToVifUpdateMap()
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
	ToVifListQuery() (string, error)
}

// ListOpts allows extensions to add additional parameters to the API.
type ListOpts struct {
	ID   string `q:"id"`
	Name string `q:"name"`
}

// ToVifListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToVifListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List is a method by which can get the detailed information of all virtual gateways.
func List(c *golangsdk.ServiceClient, opts ListOptsBuilder) ([]VirtualInterface, error) {
	url := rootURL(c)
	query, err := opts.ToVifListQuery()
	if err != nil {
		return nil, err
	}
	url += query
	var r CommonResult
	_, r.Err = c.Get(url, &r.Body, nil)

	return ExtractVirtualInterfaces(r)
}

// Delete is a method to delete an existing virtual interface.
func Delete(c *golangsdk.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, id), nil)
	return
}
