---
subcategory: "Virtual Private Cloud (VPC)"
---

# hcs_vpc_subnet

Provides a VPC subnet resource within HuaweiCloudStack.

## Example Usage

```hcl
resource "hcs_vpc" "vpc" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "hcs_vpc_subnet" "subnet" {
  name       = var.subnet_name
  cidr       = var.subnet_cidr
  gateway_ip = var.subnet_gateway_ip
  vpc_id     = hcs_vpc.vpc.id
}

 ```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies tThe region in which to create the vpc subnet. If omitted, the
  provider-level region will be used. Changing this creates a new Subnet.

* `name` - (Required, String) Specifies the subnet name. The value is a string of 1 to 64 characters that can contain
  letters, digits, underscores (_), and hyphens (-).

* `cidr` - (Required, String, ForceNew) Specifies the network segment on which the subnet resides. The value must be in
  CIDR format and within the CIDR block of the VPC. The subnet mask cannot be greater than 28. Changing this creates a
  new Subnet.

* `gateway_ip` - (Required, String, ForceNew) Specifies the gateway of the subnet. The value must be a valid IP address
  in the subnet segment. Changing this creates a new Subnet.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of the VPC to which the subnet belongs. Changing this creates
  a new Subnet.

* `description` - (Optional, String) Specifies supplementary information about the subnet. The value is a string of
  no more than 255 characters and cannot contain angle brackets (< or >).

* `ipv6_enable` - (Optional, Bool) Specifies whether the IPv6 function is enabled for the subnet. Defaults to false.

* `dhcp_enable` - (Optional, Bool) Specifies whether the DHCP function is enabled for the subnet. Defaults to true.

* `primary_dns` - (Optional, String) Specifies the IP address of DNS server 1 on the subnet. The value must be a valid
  IP address.

* `secondary_dns` - (Optional, String) Specifies the IP address of DNS server 2 on the subnet. The value must be a valid
  IP address.

* `dns_list` - (Optional, List) Specifies the DNS server address list of a subnet. This field is required if you need to
  use more than two DNS servers. This parameter value is the superset of both DNS server address 1 and DNS server
  address 2.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `status` - The status of the subnet. The value can be ACTIVE, DOWN, UNKNOWN, or ERROR.

* `ipv4_subnet_id` - The ID of the IPv4 subnet (Native OpenStack API).

* `ipv6_subnet_id` - The ID of the IPv6 subnet (Native OpenStack API).

* `ipv6_cidr` - The IPv6 subnet CIDR block.

* `ipv6_gateway` - The IPv6 subnet gateway.

## Import

Subnets can be imported using the subnet `Network ID`, e.g.

```
$ terraform import hcs_vpc_subnet.subnet_v1 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `delete` - Default is 10 minute.
