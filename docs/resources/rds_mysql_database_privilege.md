---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloudStack"
page_title: "HuaweiCloudStack: hcs_rds_mysql_database_privilege"
description: |-
  Manages RDS Mysql database privilege resource within HuaweiCloudStack.
---

# hcs_rds_mysql_database_privilege

Manages RDS Mysql database privilege resource within HuaweiCloudStack.

## Example Usage

```hcl
variable "rds_instance_id" {}
variable "db_name" {}
variable "password" {}

resource "hcs_rds_mysql_database" "test" {
  instance_id   = var.rds_instance_id
  name          = var.db_name
  character_set = "utf8"
  description   = "test database"
}

resource "hcs_rds_mysql_account" "test1" {
  instance_id = var.rds_instance_id
  name        = "test1"
  password    = var.password
}

resource "hcs_rds_mysql_account" "test2" {
  instance_id = var.rds_instance_id
  name        = "test2"
  password    = var.password
}

resource "hcs_rds_mysql_database_privilege" "test" {
  instance_id = var.rds_instance_id
  db_name     = hcs_rds_mysql_database.test.name

  users {
    name     = hcs_rds_mysql_account.test1.name
    readonly = true
  }

  users {
    name     = hcs_rds_mysql_account.test2.name
    readonly = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the RDS database privilege resource. If omitted,
  the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the RDS instance ID.

  Changing this will create a new resource.

* `db_name` - (Required, String, ForceNew) Specifies the database name.

  Changing this will create a new resource.

* `users` - (Required, List) Specifies the account that associated with the database.

  The [users](#rds_mysql_users) structure is documented below.

<a name="rds_mysql_users"></a>
The `users` block supports:

* `name` - (Required, String) Specifies the username of the database account.

* `readonly` - (Optional, Bool) Specifies the read-only permission. The value can be:
  + **true**: indicates the **read-only** permission.
  + **false**: indicates the **read and write** permission.

  The default value is **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of database privilege which is formatted `<instance_id>/<db_name>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

* `update` - Default is 30 minutes.

* `delete` - Default is 30 minutes.

## Import

RDS database privilege can be imported using the `instance id` and `db_name`, e.g.

```bash
$ terraform import hcs_rds_mysql_database_privilege.test <instance_id>/<db_name>
```
