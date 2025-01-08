---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_rds_sql_audit"
description: ""
---

# hcs_rds_sql_audit

Manages RDS SQL audit resource within HuaweiCloudStack.

-> **NOTE:** Only MySQL and PostgreSQL engines are supported.

## Example Usage

```hcl
variable "instance_id" {}

resource "hcs_rds_sql_audit" "test" {
  instance_id = var.instance_id
  keep_days   = 5

  audit_types = [
    "CREATE_USER",
    "DROP_USER",
    "DROP",
    "INSERT",
    "BEGIN/COMMIT/ROLLBACK"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the RDS instance.

  Changing this parameter will create a new resource.

* `keep_days` - (Required, Int) Specifies the number of days for storing audit logs. Value ranges from `1` to `732`.

* `audit_types` - (Optional, List) Specifies the list of audit types. This parameter applies only to **MySQL**. It is
  not supported for PostgreSQL. Value options: 
  - **CREATE_USER**, **DROP_USER**, **RENAME_USER**, **GRANT**, **REVOKE**, **ALTER_USER**, **ALTER_USER_DEFAULT_ROLE**.
  - **CREATE**, **ALTER**, **DROP**, **RENAME**, **TRUNCATE**, **REPAIR**, **OPTIMIZE**.
  - **INSERT**, **DELETE**, **UPDATE**, REPLACE, **SELECT**.
  - **BEGIN/COMMIT/ROLLBACK**, **PREPARED_STATEMENT**, **CALL_PROCEDURE**, **KILL**, **SET_OPTION**, **CHANGE_DB**,
    **UNINSTALL_PLUGIN**, **UNINSTALL_PLUGIN**, **INSTALL_PLUGIN**, **SHUTDOWN**, **SLAVE_START**, **SLAVE_STOP**,
    **LOCK_TABLES**, **UNLOCK_TABLES**, **FLUSH**, **XA**.

* `reserve_auditlogs` - (Optional, Bool) Specifies whether the historical audit logs will be reserved for some time
  when SQL audit is disabled. It is valid only when SQL audit is **disabled**. Default to **true**.
  - **true**. Indicates that historical audit logs are deleted after the audit log policy is disabled. 
  - **false**. Indicates that historical audit logs are deleted when the audit log policy is disabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The RDS SQL audit can be imported using the `id`, e.g.

```bash
$ terraform import hcs_rds_sql_audit.test <id>
```
