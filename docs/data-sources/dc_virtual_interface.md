---
subcategory: "Virtual Private Cloud (VPC)"
---

# hcs_virtual_interface

Provides details about a specific Virtual Interface.

## Example Usage

```hcl
variable "virtual_interface_name" {}

data "hcs_virtual_interface" "virtual_interface" {
  name = var.virtual_interface_name
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available Virtual Interfaces in the current region.
The given filters must match exactly one Virtual Interface whose data will be exported as attributes.

* `region` - (Optional, String) Specifies the region in which to obtain the Virtual Interface. If omitted,
  the provider-level region will be used.

* `name` - (Optional, String) Specifies an unique name for the Virtual Interface. The value is a string of no more than
  64 characters and can contain digits, letters, underscores (_), and hyphens (-).

* `id` - (Optional, String) Specifies the id of the Virtual Interface to retrieve.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:


* `description` - The description of the Virtual Interface.

* `status` - Whether the Virtual Interface connection is available. The value can be ACTIVE, DOWN, BUILD, ERROR,
  PENDING_DELETE, DELETED, APPLY, DENY, PENDING_PAY, PAID, ORDERING, ACCEPT, or REJECTED.

* `direct_connect_id` - The ID of Direct Connect bound to the Virtual Interface.

* `vgw_id` - The ID of Virtual Gateway bound to the Virtual Interface.

* `remote_ep_group` - An array of one or more remote network CIDR.

* `remote_ep_group_id` - Endpoint group ID. The endpoint group contains the CIDR of the remote network.

* `link_infos` - Interconnection information object between the L3GW and the PE.
  The [link_info object](#link_info_object).

<a name="link_info_object"></a>
the `link_infos` block supports:

* `interface_group_id` - Interface group ID.

* `hosting_id` - ID of the hosting connection.

* `local_gateway_v4_ip` - IPv4 address of the local gateway.

* `local_gateway_v6_ip` - IPv6 address of the local gateway.

* `remote_gateway_v4_ip` - IPv4 address of the remote gateway.

* `remote_gateway_v6_ip` - IPv6 address of the remote gateway.

* `vlan` - Interconnection VLAN between the L3GW TOR and the off-cloud device.

* `bgp_asn` - AS number of the BGP peer of the off-cloud device.

* `bgp_asn_dot` - AS number of the BGP peer of the off-cloud device, in X.Y format.

