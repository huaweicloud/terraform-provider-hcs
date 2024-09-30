---
subcategory: "Virtual Private Cloud (VPC)"
---

# hcs_vpc_route_table

Manages a VPC custom route table resource within hcs.

## Example Usage

### Basic Custom Route Table

```hcl
variable "vpc_id" {}
variable "peer_vpc_id" {}

resource "hcs_vpc_route_table" "demo" {
  name        = "demo"
  vpc_id      = var.vpc_id
  description = "a custom route table demo"

  route {
    destination = "172.16.0.0/16"
    type        = "peering"
    nexthop     = var.peer_vpc_id
  }
}
```

### Associating Subnets with a Route Table

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "peer_vpc_id" {}

data "hcs_vpc_subnet" "subnet" {
  id = var.subnet_id
}

resource "hcs_vpc_route_table" "demo" {
  name    = "demo"
  vpc_id  = var.vpc_id
  subnets = [data.hcs_vpc_subnet.subnet.ipv4_subnet_id]

  route {
    destination = "172.16.0.0/16"
    type        = "peering"
    nexthop     = var.peer_vpc_id
  }
  route {
    destination = "192.168.100.0/24"
    type        = "vip"
    nexthop     = "192.168.10.200"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the vpc route table.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID for which a route table is to be added.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the route table name. The value is a string of no more than
  `64` characters that can contain letters, digits, underscores (_), hyphens (-), and periods (.).

* `description` - (Optional, String) Specifies the supplementary information about the route table.
  The value is a string of no more than `255` characters and cannot contain angle brackets (< or >).

* `subnets` - (Optional, List) Specifies an array of one or more neutron subnets associating with the route table.

  -> **NOTE:** The custom route table associated with a subnet affects only the outbound traffic.
  The default route table determines the inbound traffic.

* `route` - (Optional, List) Specifies the route object list. The [route object](#route_object)
  is documented below.

<a name="route_object"></a>
The `route` block supports:

* `destination` - (Required, String) Specifies the destination address in the CIDR notation format,
  for example, 192.168.200.0/24. The destination of each route must be unique and cannot overlap
  with any subnet in the VPC.

* `type` - (Required, String) Specifies the route type. Currently, the value can be:
  **eni**, **subeni**, **vip**, **nat**, **peering**, **vpn**, **dc**, **cc** and **externalip**.

* `nexthop` - (Required, String) Specifies the next hop.
  + If the route type is **eni**, the value is the NIC or extension NIC of an ECS in the VPC. 
  + If the route type is **subeni**, the value is the supplementary NIC of a NIC in the VPC.
  + If the route type is **vip**, the value is a VIP port ID.
  + If the route type is **nat**, the value is a VPN gateway ID.
  + If the route type is **peering**, the value is a peer VPC ID.
  + If the route type is **vpn**, the value is a VPN gateway ID.
  + If the route type is **dc**, the value is a Direct Connect gateway ID.
  + If the route type is **cc**, the value is a Cloud Connection ID.
  + If the route type is **externalip**, the value is an external IP address.

* `description` - (Optional, String) Specifies the supplementary information about the route.
  The value is a string of no more than 255 characters and cannot contain angle brackets (< or >).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

vpc route tables can be imported using the `id`, e.g.

```bash
$ terraform import hcs_vpc_route_table.demo e1b3208a-544b-42a7-84e6-5d70371dd982
```
