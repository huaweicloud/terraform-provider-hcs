---
subcategory: "Virtual Private Cloud (VPC)"
---

# hcs_direct_connect

Manages a Direct Connect resource within hcs.

## Example Usage

```hcl
variable "hosting_direct_connect_id" {}
variable "peer_location" {}

resource "hcs_direct_connect" "demo" {
    name = "demo"
    description = "a direct connect demo"
    hosting_id = var.hosting_direct_connect_id
    type = "hosted"
    peer_location = var.peer_location
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the Direct Connect. If omitted, the
  provider-level region will be used. Changing this creates a new Direct Connect resource.

* `name` - (Optional, String) Specifies the name of the Direct Connect. The value can contain 1 to 64 characters.

* `description` - (Optional, String) Specifies the description of the Direct Connect.
  The value is a string of no more than `128` characters and cannot contain angle brackets (< or >).

* `hosting_id` - (Required, String, ForceNew) Specifies the ID of Direct Connect access point bound to the Direct Connect.
  Changing this creates a new Direct Connect resource.

* `type` - (Required, String, ForceNew) Specifies the type of the Direct Connect. the value can be "hosted".
  Changing this creates a new Direct Connect resource.

* `peer_location` - (Optional, String) Specifies the user network location of the Direct Connect.

* `tenancy` - (Optional, String) Specifies the lease period of the Direct Connect.

## Attribute Reference

* `id` - The resource ID in UUID format.

* `status` - Whether the Direct Connect is available. The value can be ACTIVE, DOWN, BUILD, ERROR,
  PENDING_DELETE, DELETED, APPLY, DENY, PENDING_PAY, PAID, ORDERING, ACCEPT, or REJECTED.

* `dc_provider` - Direct Connect type. The value can be "ce" or "vpc-gw".

* `group` - Egress of the Direct Connect.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

Direct Connects can be imported using the `id`, e.g.

```
$ terraform import hcs_direct_connect.demo 7117d38e-4c8f-4624-a505-bd96b97d024c
```
