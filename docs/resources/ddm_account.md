---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_ddm_account"
description: |-
  Use this resource to set a DDM instance read strategy within HuaweiCloudStack.
---

# hcs_ddm_account

Manages a DDM account resource within HuaweiCloudStack.

## Example Usage

```hcl
variable "ddm_instance_id" {}
variable "name" {}
variable "password" {}
variable "schema_name" {}

resource "hcs_ddm_account" "test"{
  instance_id = var.ddm_instance_id
  name        = var.name
  password    = var.password

  permissions = [
    "CREATE",
    "SELECT"
  ]

  schemas {
   name = var.schema_name
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of a DDM instance.
  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the DDM account.
  An account name starts with a letter, consists of 1 to 32 characters, and can contain only letters,
  digits, and underscores (_).
  Changing this parameter will create a new resource.

* `password` - (Required, String) Specifies the DDM account password.

* `permissions` - (Required, List) Specifies the basic permissions of the DDM account.
  The valid valuse are as follows:
  + **CREATE**
  + **DROP**
  + **ALTER**
  + **INDEX**
  + **INSERT**
  + **DELETE**
  + **UPDATE**
  + **SELECT**

* `description` - (Optional, String) Specifies the description of the DDM account.

* `databases` - (Optional, List) Specifies the databases that associated with the account.
  The [databases](#ddm_account_databases_arg) structure is documented below.

<a name="ddm_account_databases_arg"></a>
The `databases` block supports:

* `name` - (Optional, String) Specifies the name of the associated schema.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Indicates the status of the DDM account.

* `password_last_changed` - The password change time of the DDM account.

* `created` - The creation time of the DDM account.

* `databases` - (Optional, List) Specifies the databases that associated with the account.
  The [databases](#ddm_account_databases_attr) structure is documented below.

<a name="ddm_account_databases_attr"></a>
The `databases` block supports:

* `description` - (Optional, String) Specifies the schema description.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

* `update` - Default is 10 minutes.

* `delete` - Default is 10 minutes.

## Import

The DDM account can be imported using the instance ID and account name separated by a slash, e.g.:

```
$ terraform import hcs_ddm_account.test 0a8f1c6baa124e99853719d9257324dfin09/account_name
```
