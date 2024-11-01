---
subcategory: "Virtual Private Cloud (VPC)"
---

# hcs_vpc_peering

Provides a resource to manage a VPC Peering resource.

## Example Usage

 ```hcl
 var peer_conn_name {}
 var vpc_id {}
 var accepter_vpc_id {}
 var peer_region {}
 
resource "hcs_vpc_peering" "peering" {
  name        = var.peer_conn_name
  vpc_id      = var.vpc_id
  peer_vpc_id = var.accepter_vpc_id
  peer_region = var.peer_region
}
 ```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the VPC peering connection. If omitted, the
  provider-level region will be used. Changing this creates a new VPC peering connection resource.

* `name` - (Required, String) Specifies the name of the VPC peering . The value can contain 1 to 64
  characters.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of a VPC involved in a VPC peering connection. Changing this
  creates a new VPC peering.

* `peer_vpc_id` - (Required, String, ForceNew) Specifies the VPC ID of the peering. Changing this creates a new
  VPC peering.

* `peer_project_id` - (Optional, String, ForceNew) Specifies the project ID of the accepter project. Changing this creates
  a new VPC peering connection.

* `peer_region` - (Optional, String, ForceNew) Specifies name of the project to which the peer VPC belongs. This parameter is 
  mandatory when the two VPCs are not in the same project.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The VPC peering ID.

* `status` - The VPC peering status. The value can be active.

## Timeouts

## Import

VPC Peering can be imported using the `id`, e.g.

```
$ terraform import hcs_vpc_peering.test 2c7f39f3-702b-48d1-940c-b50384177ee1
```