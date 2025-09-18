---
subcategory: "Virtual Private Cloud (VPC)"
---

# hcs_virtual_gateway

Provides details about a specific Virtual Gateway.

## Example Usage

```hcl
variable "virtual_gateway_name" {}

data "hcs_virtual_gateway" "virtual_gateway" {
  name = var.virtual_gateway_name
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available Virtual Gateways in the current region.
The given filters must match exactly one Virtual Gateway whose data will be exported as attributes.

* `region` - (Optional, String) Specifies the region in which to obtain the Virtual Gateway. If omitted,
  the provider-level region will be used.

* `name` - (Optional, String) Specifies an unique name for the Virtual Gateway. The value is a string of no more than
  64 characters and can contain digits, letters, underscores (_), and hyphens (-).

* `id` - (Optional, String) Specifies the id of the Virtual Gateway to retrieve.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` - Whether the Virtual Gateway connection is available. The value can be ACTIVE, DOWN, BUILD, ERROR,
  PENDING_CREATE, PENDING_UPDATE or PENDING_DELETE.

* `description` - (Optional, String) Specifies the description of the Virtual Gateway.

* `vpc_group` - An array of one or more vpc_group objects binding with the Virtual Gateway.
  The [vpc_group object](#vpc_group_object).

<a name="vpc_group_object"></a>
The `vpc_group` block supports:

* `vpc_id` - The VPC ID.

* `local_ep_group` - An array of one or more local network VPC CIDR.
