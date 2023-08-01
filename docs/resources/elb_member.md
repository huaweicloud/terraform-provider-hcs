---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
---

# hcs_elb_member

Manages an ELB member resource within HCS.

## Example Usage

```hcl
variable "elb_pool_id" {}
variable "ipv4_subnet_id" {}

resource "hcs_elb_member" "member_1" {
  address       = "192.168.199.23"
  protocol_port = 8080
  pool_id       = var.elb_pool_id
  subnet_id     = var.ipv4_subnet_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the ELB member resource. If omitted, the the
  provider-level region will be used. Changing this creates a new member.

* `pool_id` - (Required, String, ForceNew) The id of the pool that this member will be assigned to.

* `subnet_id` - (Optional, String, ForceNew) The **IPv4 or IPv6 subnet ID** of the subnet in which to access the member.
  + The IPv4 or IPv6 subnet must be in the same VPC as the subnet of the load balancer.
  + If this parameter is not specified, **cross-VPC backend** has been enabled for the load balancer.
    In this case, cross-VPC backend servers must use private IPv4 addresses,
    and the protocol of the backend server group must be TCP, HTTP, or HTTPS.

* `name` - (Optional, String) Human-readable name for the member.

* `address` - (Required, String, ForceNew) The IP address of the member to receive traffic from the load balancer.
  Changing this creates a new member.

* `protocol_port` - (Required, Int, ForceNew) The port on which to listen for client traffic. Changing this creates a
  new member.

* `weight` - (Optional, Int)  A positive integer value that indicates the relative portion of traffic that this member
  should receive from the pool. For example, a member with a weight of 10 receives five times as much traffic as a
  member with a weight of 2.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the member.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `update` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

ELB member can be imported using the pool ID and member ID separated by a slash, e.g.

```
$ terraform import hcs_elb_member.member_1 e0bd694a-abbe-450e-b329-0931fd1cc5eb/4086b0c9-b18c-4d1c-b6b8-4c56c3ad2a9e
```
