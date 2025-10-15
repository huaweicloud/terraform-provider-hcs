---
subcategory: "Virtual Data Center (VDC)"
---

# hcs_vdc_group

Use this data source to get details of the specified user group

## Example Usage

```hcl
variable "group_name" {}
variable "vdc_id" {}

data "hcs_vdc_group" "group_1" {
  vdc_id = var.vdc_id
  name   = var.group_name
}
```

## Argument Reference

The following arguments are supported:

* `vdc_id` - (Required, String) VDC ID.

* `name` - (Required, String) User group name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - User group ID.

* `domain_id` - Tenant ID.

* `description` - User group description.

* `create_at` - User group creation time.
