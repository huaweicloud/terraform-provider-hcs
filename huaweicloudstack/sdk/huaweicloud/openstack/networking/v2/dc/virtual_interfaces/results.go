package virtual_interfaces

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

type CommonResult struct {
	golangsdk.Result
}

type DeleteResult struct {
	golangsdk.ErrResult
}

// Extract is a function that accepts a result and extracts a VirtualInterface resource.
func (r CommonResult) Extract() (*VirtualInterface, error) {
	var s VirtualInterface
	err := r.ExtractInto(&s)
	return &s, err
}

func (r CommonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "virtual_interface")
}

// ExtractVirtualInterfaces is a function that accepts a result and extracts a slice of VirtualInterfaces structs.
func ExtractVirtualInterfaces(r CommonResult) ([]VirtualInterface, error) {
	var s []VirtualInterface
	err := ExtractVirtualInterfacesInto(r, &s)
	return s, err
}

func ExtractVirtualInterfacesInto(r CommonResult, v interface{}) error {
	return r.Result.ExtractIntoSlicePtr(v, "virtual_interfaces")
}

// VirtualInterface is a struct that represents the detail of the virtual interface.
type VirtualInterface struct {
	// UUID for the virtual interface.
	ID string `json:"id"`
	// Human-readable name for the virtual interface. Might not be unique.
	Name string `json:"name"`

	// Indicates whether virtual interface is currently operational. Possible values include
	// 'ACTIVE', 'DOWN', 'BUILD', 'ERROR', 'PENDING_CREATE', 'PENDING_UPDATE', 'PENDING_DELETE',
	// 'DELETED', 'AUTHORIZATION', 'REJECTED'.
	Status string `json:"status"`

	// User-defined description of the virtual interface.
	Description string `json:"description"`

	// UUID of the direct connect bound with the virtual interface.
	DirectConnectId string `json:"direct_connect_id"`

	// UUID of the virtual gateway bound with the virtual interface.
	VgwId string `json:"vgw_id"`

	// UUID of the dc endpoint group bound by the virtual interface.
	RemoteEpGroupId string `json:"remote_ep_group_id"`

	// Link infos bound by the virtual interface.
	LinkInfos []LinkInfo `json:"link_infos"`
}

// LinkInfo is a struct that represents the link info of the virtual interface.
type LinkInfo struct {
	// UUID of the interface group to be bound.
	InterfaceGroupId string `json:"interface_group_id"`
	// UUID of the hosting direct connect.
	HostingId string `json:"hosting_id"`
	// Ipv4 local gateway ip.
	LocalGatewayV4Ip string `json:"local_gateway_v4_ip"`
	// Ipv6 local gateway ip.
	LocalGatewayV6Ip string `json:"local_gateway_v6_ip"`
	// Ipv4 remote gateway ip.
	RemoteGatewayV4Ip string `json:"remote_gateway_v4_ip"`
	// Ipv6 remote gateway ip.
	RemoteGatewayV6Ip string `json:"remote_gateway_v6_ip"`
	// VLAN to be used.
	Vlan int `json:"vlan"`
	// BGP peer as, numeric.
	BgpAsn int `json:"bgp_asn"`
	// BGP peer as, dotted.
	BgpAsnDot string `json:"bgp_asn_dot"`
}
