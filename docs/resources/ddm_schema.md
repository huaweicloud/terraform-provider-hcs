---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_ddm_schema"
description: |-
  Manages a DDM schema resource within HuaweiCloudStack.
---

# hcs_ddm_schema

Manages a DDM schema resource within HuaweiCloudStack.

## Example Usage

```hcl
variable "ddm_instance_id" {}
variable "rds_instance_id" {}
variable "password" {}

resource "hcs_ddm_schema" "test"{
  instance_id  = var.ddm_instance_id
  name         = "test_schema"
  shard_mode   = "single"
  shard_number = 8
  shard_unit   = 8

  data_nodes {
    id             = var.rds_instance_id
    admin_user     = "root"
    admin_password = var.password
  }

  delete_rds_data = "true"
  
  lifecycle {
    ignore_changes = [
      data_nodes,
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of a DDM instance.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the DDM schema.
  An instance name starts with a letter, consists of `2` to `48` characters, and can contain only lowercase letters,
  digits, and underscores (_). Cannot contain keywords information_schema, mysql, performance_schema, or sys.

  Changing this parameter will create a new resource.

* `shard_mode` - (Required, String, ForceNew) Specifies the sharding mode of the schema.  
  The valid values are as follows:
  + **cluster**: Indicates that the schema is in sharded mode.
  + **single**: Indicates that the schema is in non-sharded mode.

  Changing this parameter will create a new resource.

* `shard_number` - (Required, Int, ForceNew) Specifies the number of shards in the same working mode.
  + if `shard_unit` is not empty, `shard_number` is the product of the value of `shard_unit` and
    the number of associated RDS instances.
  + if `shard_unit` is empty, `shard_number` must be greater than the number of associated RDS instances and
    less than or equal to the number of associated RDS instances multiplied by `64`.

  Changing this parameter will create a new resource.

* `shard_unit` - (Optional, Int, ForceNew) Specifies the number of logical database shards on RDS instance.

  Changing this parameter will create a new resource.

  ->**Note** In HCS **8.3.1** and **8.5.0**, the `shard_unit` is a **required** parameter.
  + if in sharded mode, the value of `shard_unit` is fixed at `1`.
  + if in non-sharded mode, the value of `shard_unit` is `8` or `16`.

  ->**Note** In HCS **8.5.1** and **8.6.0**, the `shard_unit` is an **optional** parameter.
  + if in sharded mode, the value of `shard_unit` is fixed at `1`.
  + if in non-sharded mode, the value of `shard_unit` is range from `1` to `64`.

* `data_nodes` - (Required, List, ForceNew) Specifies the RDS instances associated with the schema.
  The [data_nodes](#ddm_schema_data_nodes_arg) structure is documented below.

  Changing this parameter will create a new resource.

* `delete_rds_data` - (Optional, String) Specifies whether data stored on the associated DB instances is deleted.

<a name="ddm_schema_data_nodes_arg"></a>
The `data_nodes` block supports:

* `id` - (Required, String) Specifies the ID of the RDS instance associated with the schema.

* `admin_user` - (Required, String) Specifies the username for logging in to the associated RDS instance.

* `admin_password` - (Required, String) Specifies the password for logging in to the associated RDS instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Indicates the schema status.

* `shards` - Indicates the sharding information of the schema.
  The [shards](#ddm_schema_shards_attr) structure is documented below.

* `data_nodes` - Indicates the RDS instances associated with the schema.
  The [data_nodes](#ddm_schema_data_nodes_attr) structure is documented below.

* `data_vips` - Indicates the IP address and port number for connecting to the schema.

<a name="ddm_schema_shards_attr"></a>
The `shards` block supports:

* `id` - Indicates the ID of the RDS instance where the shard is located.

* `id_name` - Indicates the name of the physical database.

* `name` - Indicates the shard name.

* `status` - Indicates the shard status.

* `db_slot` - Indicates the number of shards.

* `created` - The creation time of the database.

* `updated` - The update time of the database.

<a name="ddm_schema_data_nodes_attr"></a>
The `data_nodes` block supports:

* `id` - Indicates the ID of the RDS instance associated with the schema.

* `name` - Indicates the name of the associated RDS instance.

* `status` - Indicates the status of the associated RDS instance.

* `error_msg` - The error message of the RDS instance. This parameter is not returned if
  no abnormal information is found.
 
## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.

* `delete` - Default is 30 minutes.

## Import

The DDM schema can be imported using the `<instance_id>/<schema_name>`, e.g.

```
$ terraform import hcs_ddm_schema.test 80e373f9-872e-4046-aae9-ccd9ddc55511/schema_name
```
