---
subcategory: "Virtual Private Cloud (VPC)"
---

# hcs_virtual_gateway

Manages a Virtual Gateway resource within hcs.

## Example Usage

```hcl
variable "vpc_id" {}

locals {
  vpc_group = [
    {
      vpc_id = var.vpc_id,
      local_ep_group = ["192.168.1.0/24","fc00:1::/64"]
    },
  ]
}

resource "hcs_virtual_gateway" "demo" {
    name = "demo"
    description = "a Virtual Gateway demo"
    dynamic "vpc_group" {
        for_each = local.vpc_group
        content {
            vpc_id = vpc_group.value.vpc_id
            local_ep_group = vpc_group.value.local_ep_group
        }
    }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the Virtual Gateway. If omitted, the
  provider-level region will be used. Changing this creates a new Virtual Gateway resource.

* `name` - (Optional, String) Specifies the name of the Virtual Gateway. The value can contain 1 to 64 characters.

* `description` - (Optional, String) Specifies the description of the Virtual Gateway.
  The value is a string of no more than `128` characters and cannot contain angle brackets (< or >).

* `vpc_group` - (Required, Set) Specifies a set of one or more vpc_group objects binding with the Virtual Gateway.
  The [vpc_group object](#vpc_group_object).

<a name="vpc_group_object"></a>
The `vpc_group` block supports:

* `vpc_id` - (Required, String) Specifies the VPC ID.

* `local_ep_group` - (Required, String) Specifies a set of one or more local network VPC CIDR.


## Attribute Reference

* `id` - The resource ID in UUID format.

* `status` - Whether the Virtual Gateway connection is available. The value can be ACTIVE, DOWN, BUILD, ERROR,
  PENDING_CREATE, PENDING_UPDATE or PENDING_DELETE.

## Timeouts

## Import

Virtual Gateways can be imported using the `id`, e.g.

```
$ terraform import hcs_virtual_gateway.demo 7117d38e-4c8f-4624-a505-bd96b97d024c
```
