---
subcategory: "VDC"
---

# hcs_vdc_group_membership

Manages VDC user group membership within Huawei Cloud Stack.

> [!NOTE]
>
> Supported from ManageOne version 8.5.1 onwards.

## Example Usage

```hcl
variable "user_password" {}
variable "vdc_id" {
  default = "a18c2ce0-5379-4b34-8a12-eee47f5cfa89"
}
resource "hcs_vdc_user" "user01" {
  vdc_id = var.vdc_id
  name = "Username1"
  password = var.user_password
}

resource "hcs_vdc_user" "user02" {
  vdc_id = var.vdc_id
  name = "Username2"
  password = var.user_password
}

resource "hcs_vdc_group" "group01" {
  vdc_id       = var.vdc_id
  name         = "Usergroup1"
}

resource "hcs_vdc_group_membership" "group_membership_1" {
  group = hcs_vdc_group.group01.id
  users = [hcs_vdc_user.user01.id, hcs_vdc_user.user02.id]
}

```

## Argument Reference

The following arguments are supported:

* `group` - (Required, String, ForceNew) User group ID.

* `users` - (Required, Set) User ID list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - User group ID.

## Import

VDC user group membership can be imported using the `id`, e.g.

```
$ terraform import hcs_vdc_group_membership.vdc_group_membership1 3b002f5e4aae407082630a00d2ac0f40
```
