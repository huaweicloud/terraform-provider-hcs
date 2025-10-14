---
subcategory: "Virtual Data Center (VDC)"
---

# hcs_vdc_group

Manages a VDC user group within Huawei Cloud Stack.

## Example Usage

```hcl
variable "vdc_id" {}

variable "group_name" {}

resource "hcs_vdc_group" "group01" {
  vdc_id      = var.vdc_id
  name        = var.group_name
  description = "Description"
}
```

## Argument Reference

The following arguments are supported:

* `vdc_id` - (Required, String) VDC ID. The value can contain 1 to 36 characters, including only lowercase letters,
  digits, and hyphens (-). Once set, the value cannot be modified.

* `name` - (Required, String) User group name. Enter 1 to 64 characters. Only letters, digits, hyphens (-), and
  underscores (_) are allowed. The value cannot start with a digit or be `admin`, `power_user`, or `guest`.

* `description` - (Optional, String) Description. The value cannot contain the following characters: >< The value can
  contain 0 to 255 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - User group ID.

## Import

Groups can be imported using the id, e.g.

```bash
$ terraform import hcs_vdc_group.group02 1ff4536fb0a44faba80450f9da0bf47a
```
