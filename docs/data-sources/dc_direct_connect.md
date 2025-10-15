---
subcategory: "Virtual Private Cloud (VPC)"
---

# hcs_direct_connect

Provides details about a specific Direct Connect.

## Example Usage

```hcl
variable "direct_connect_name" {}

data "hcs_direct_connect" "direct_connect" {
  name = var.direct_connect_name
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available Direct Connects in the current region.
The given filters must match exactly one Direct Connect whose data will be exported as attributes.

* `region` - (Optional, String) Specifies the region in which to obtain the Direct Connect. If omitted,
  the provider-level region will be used.

* `name` - (Optional, String) Specifies an unique name for the Direct Connect. The value is a string of no more than
  64 characters and can contain digits, letters, underscores (_), and hyphens (-).

* `id` - (Optional, String) Specifies the id of the Direct Connect to retrieve.

* `dc_provider` - (Optional, String) Specifies the dc_provider desired Direct Connect. The value can be "ce".

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `description` - The description of the Direct Connect.

* `status` - Whether the Direct Connect is available. The value can be ACTIVE, DOWN, BUILD, ERROR,
  PENDING_DELETE, DELETED, APPLY, DENY, PENDING_PAY, PAID, ORDERING, ACCEPT, or REJECTED.

* `hosting_id` - The ID of Direct Connect access point bound to the Direct Connect.

* `type` - The type of the Direct Connect. the value can be "hosted".

* `peer_location` - The user network location of the Direct Connect.

* `group` - Egress of the Direct Connect.
