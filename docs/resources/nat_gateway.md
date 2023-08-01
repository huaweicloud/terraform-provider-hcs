---
subcategory: "NAT Gateway (NAT)"
---

# hcs_nat_gateway

Manages a gateway resource of the **public** NAT within HuaweiCloudStack(hcs).

## Example Usage

```hcl
variable "gateway_name" {}
variable "vpc_id" {}
variable "network_id" {}

resource "hcs_nat_gateway" "test" {
  name        = var.gateway_name
  description = "test for terraform"
  spec        = "3"
  vpc_id      = var.vpc_id
  subnet_id   = var.network_id
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of the VPC to which the NAT gateway belongs.  
  Changing this will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the subnet ID of the downstream interface (the next hop of the
  DVR) of the NAT gateway.  
  Changing this will create a new resource.

* `name` - (Required, String) Specifies the NAT gateway name.  
  The valid length is limited from `1` to `64`, only letters, digits, hyphens (-) and underscores (_) are allowed.

* `spec` - (Required, String) Specifies the specification of the NAT gateway. The valid values are as follows:
  + **1**: Small type, which supports up to `10,000` SNAT connections.
  + **2**: Medium type, which supports up to `50,000` SNAT connections.
  + **3**: Large type, which supports up to `200,000` SNAT connections.
  + **4**: Extra-large type, which supports up to `1,000,000` SNAT connections.

* `description` - (Optional, String) Specifies the description of the NAT gateway, which contain maximum of `255`
  characters, and angle brackets (<) and (>) are not allowed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `status` - The current status of the NAT gateway.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

NAT gateways can be imported using their `id`, e.g.

```bash
$ terraform import hcs_nat_gateway.test d126fb87-43ce-4867-a2ff-cf34af3765d9
```
