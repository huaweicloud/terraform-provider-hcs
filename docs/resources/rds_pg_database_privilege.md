---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloudStack"
page_title: "HuaweiCloudStack: hcs_rds_pg_database_privilege"
description: |-
  Manages an RDS PostgreSQL database privilege resource within HuaweiCloudStack.
---

# hcs_rds_pg_database_privilege

-> **NOTE:** This resource can only be used in HCS **8.5.0** and **later** version.

-> This resource is a one-time action resource for batch grant privileges to accounts. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.

Manages an RDS PostgreSQL database privilege resource within HuaweiCloudStack.

## Example Usage

```hcl
variable "instance_id" {}
variable "db_name" {}
variable "user_name_1" {}
variable "user_name_2" {}
variable "schema_name_1" {}
variable "schema_name_2" {}

resource "hcs_rds_pg_database_privilege" "test" {
  instance_id = var.instance_id
  db_name     = var.db_name

  users {
    name        = var.user_name_1
    schema_name = var.schema_name_1
    readonly    = true
  }

  users {
    name        = var.user_name_2
    schema_name = var.schema_name_2
    readonly    = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the RDS database privilege resource. If omitted,
  the provider-level region will be used.

  Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the RDS instance ID.

* `db_name` - (Required, String, NonUpdatable) Specifies the database name.

* `users` - (Required, List) Specifies the account that associated with the database. A single request supports a
  maximum of **50** elements.

  The [users](#rds_pg_database_privilege_users) structure is documented below.

<a name="rds_pg_database_privilege_users"></a>
The `users` block supports:

* `name` - (Required, String) Specifies the username of the database account.

* `readonly` - (Required, Bool) Specifies the read-only permission.  
  The valid values are as follows:
  + **true**: indicates the read-only permission.
  + **false**: indicates the read and write permission.

* `schema_name` - (Required, String) Specifies the name of the schema.
  The schema name must be between `1` and `63` characters long, consisting of letters, numbers, or underscores. It
  cannot contain any other special characters, start with **pg** or a number, or have the same name as the RDS for
  PostgreSQL template library. Additionally, the schema name must already exist.

  -> **Note**: The `schema_name` parameter must be provided. Otherwise, the API will report an error.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of database privilege which is formatted `<instance_id>/<db_name>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

* `update` - Default is 30 minutes.

* `delete` - Default is 30 minutes.
