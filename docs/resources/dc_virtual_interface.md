---
subcategory: "Virtual Private Cloud (VPC)"
---

# hcs_virtual_interface

Manages a Virtual Interface resource within hcs.

## Example Usage

```hcl
variable "direct_connect_id" {}
variable "virtual_gateway_id" {}
variable "hosting_direct_connect_id" {}
variable "interface_group_id1" {}
variable "interface_group_id2" {}

locals {
  link_infos = [
    {
        interface_group_id = var.interface_group_id1
        hosting_id = var.hosting_direct_connect_id
        local_gateway_v4_ip = "192.168.100.1/30"
        local_gateway_v6_ip = "fc00:1::1/64"
        remote_gateway_v4_ip = "192.168.100.2/30"
        remote_gateway_v6_ip = "fc00:1::2/64"
        vlan = 1111
        bgp_asn = 65111
    },
    {
        interface_group_id = var.interface_group_id2
        hosting_id = var.hosting_direct_connect_id
        local_gateway_v4_ip = "192.168.100.5/30"
        local_gateway_v6_ip = "fc00:2::1/64"
        remote_gateway_v4_ip = "192.168.100.6/30"
        remote_gateway_v6_ip = "fc00:2::2/64"
        vlan = 1111
        bgp_asn = 65111
    },
  ]
}

resource "hcs_virtual_interface" "demo" {
    name = "demo"
    description = "a direct connect demo"
    direct_connect_id = var.direct_connect_id
    vgw_id = var.virtual_gateway_id
    remote_ep_group = ["192.168.0.0/24","fc00::/64"]
    dynamic "link_infos" {
        for_each = local.link_infos
        content {
            interface_group_id = link_infos.value.interface_group_id
            hosting_id = link_infos.value.hosting_id
            local_gateway_v4_ip = link_infos.value.local_gateway_v4_ip
            local_gateway_v6_ip = link_infos.value.local_gateway_v6_ip
            remote_gateway_v4_ip = link_infos.value.remote_gateway_v4_ip
            remote_gateway_v6_ip = link_infos.value.remote_gateway_v6_ip
            vlan = link_infos.value.vlan
            bgp_asn = link_infos.value.bgp_asn
        }
    }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the Virtual Interface. If omitted, the
  provider-level region will be used. Changing this creates a new Virtual Interface resource.

* `name` - (Optional, String) Specifies the name of the Virtual Interface. The value can contain 1 to 64 characters.

* `description` - (Optional, String) Specifies the description of the Virtual Interface.
  The value is a string of no more than `128` characters and cannot contain angle brackets (< or >).

* `direct_connect_id` - (Required, String, ForceNew) Specifies the ID of Direct Connect bound to the Virtual Interface.
  Changing this creates a new Virtual Interface resource.

* `vgw_id` - (Required, String, ForceNew) Specifies the ID of Virtual Gateway bound to the Virtual Interface.
  Changing this creates a new Virtual Interface resource.

* `remote_ep_group` - (Optional, Set) Specifies an Array of one or more remote network CIDR.

* `link_infos` - (Optional, Set, ForceNew) Specifies Interconnection information object between the L3GW and the PE.
  If automatic allocation of interconnection information is configured for the L3GW, this field is not required.
  If this parameter is set, hosting_id, vlan, interface_group_id, remote_gateway_v4_ip and local_gateway_v4_ip are mandatory.
  Changing this creates a new Virtual Interface resource.
  The [link_info object](#link_info_object).

<a name="link_info_object"></a>
the `link_infos` block supports:

* `interface_group_id` - (Required, String, ForceNew) Interface group ID. If link_infos is set, this parameter is mandatory
  and the corresponding preconfiguration must exist. Changing this creates a new Virtual Interface resource.

* `hosting_id` - (Required, String, ForceNew) ID of the hosting connection. If link_infos is set, this parameter is mandatory
  and the corresponding preconfiguration must exist. Changing this creates a new Virtual Interface resource.

* `local_gateway_v4_ip` - (Required, String, ForceNew) IPv4 address of the local gateway. If link_infos is set,
  this parameter is mandatory and the corresponding preconfiguration must exist. Changing this creates a new
  Virtual Interface resource.

* `local_gateway_v6_ip` - (Optional, String, ForceNew) IPv6 address of the local gateway. If this parameter is set,
  the value must be the same as that in the corresponding preconfiguration. Otherwise, an error is reported.
  Changing this creates a new Virtual Interface resource.

* `remote_gateway_v4_ip` - (Required, String, ForceNew) IPv4 address of the remote gateway. If link_infos is set,
  this parameter is mandatory and the corresponding preconfiguration must exist. Changing this creates a new
  Virtual Interface resource.

* `remote_gateway_v6_ip` - (Optional, String, ForceNew) IPv6 address of the remote gateway. If this parameter is set,
  the value must be the same as that in the corresponding preconfiguration. Otherwise, an error is reported.
  Changing this creates a new Virtual Interface resource.

* `vlan` - (Required, Int, ForceNew) Interconnection VLAN between the L3GW TOR and the off-cloud device. If link_infos is set,
  this parameter is mandatory and the corresponding preconfiguration must exist. The value ranges from 1 to 4063.
  Changing this creates a new Virtual Interface resource.

* `bgp_asn` - (Optional, Int, ForceNew) AS number of the BGP peer of the off-cloud device. The value is an integer.
  The value can range from 0 to 4294967295, where 1 to 4294967295 indicate valid AS numbers and 0 indicates that
  the AS number is not specified. This parameter and bgp_asn_dot cannot be specified at the same time.
  If you set this parameter, the value must be the same as that in the corresponding preconfiguration.
  Otherwise, an error will be reported. Changing this creates a new Virtual Interface resource.

* `bgp_asn_dot` - (Optional, String, ForceNew) AS number of the BGP peer of the off-cloud device, in X.Y format.
  The value of X can range from 1 to 65535, and that of Y can range from 0 to 65535.
  The value null indicates that the parameter is not specified. If you set this parameter,
  the value must be the same as that in the corresponding preconfiguration. Otherwise, an error will be reported.
  Changing this creates a new Virtual Interface resource.

## Attribute Reference

* `id` - The resource ID in UUID format.

* `status` - Whether the Virtual Interface connection is available. The value can be ACTIVE, DOWN, BUILD, ERROR,
  PENDING_DELETE, DELETED, APPLY, DENY, PENDING_PAY, PAID, ORDERING, ACCEPT, or REJECTED.

* `remote_ep_group_id` - Endpoint group ID. The endpoint group contains the CIDR of the remote network.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `update` - Default is 20 minutes.
* `delete` - Default is 20 minutes.

## Import

Virtual Interfaces can be imported using the `id`, e.g.

```
$ terraform import hcs_direct_connect.demo 7117d38e-4c8f-4624-a505-bd96b97d024c
```
