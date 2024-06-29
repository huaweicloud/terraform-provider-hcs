---
subcategory: "Virtual Private Cloud (VPC)"
---

# hcs_vpc_peering_route

Provides a resource to manage a VPC Peering Route resource.

## Example Usage

 ```hcl
resource "hcs_vpc_peering_route" "peering" {
  peering_id  = var.vpc_peering_id
  vpc_id      = var.vpc_id
  route {
      nexthop     = var.vpc_peering_id1
      destination = "192.168.0.0/24"
  }
  route {
      nexthop     = var.vpc_peering_id2
      destination = "192.168.1.0/24"
  }
 ```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the VPC peering connection. If omitted, the
  provider-level region will be used. Changing this creates a new VPC peering connection resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of a VPC involved in a VPC peering.

* `peering_id` - (Required, String, ForceNew) Specifies the VPC Peering ID.

* `route` - (Optional, Set) Specifies the VPC Peering route.

## Timeouts

## Import

VPC Peering Route can be imported using the `{peering_id}/{vpc_id}`, e.g.

```
$ terraform import hcs_vpc_peering_route.test ff8080829058bd2401905d36e356003e/2c7f39f3-702b-48d1-940c-b50384177ee1
```