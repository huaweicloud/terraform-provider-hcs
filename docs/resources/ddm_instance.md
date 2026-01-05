---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_ddm_instance"
description: |-
  Manages DDM instance resource within HuaweiCloudStack
---

# hcs_ddm_instance

Manages DDM instance resource within HuaweiCloudStack.

## Example Usage

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "sec_group_id" {}
variable "eps_id" {}
variable "availability_zone" {}
variable "engine_id" {}
variable "flavor_id" {}

resource "hcs_ddm_instance" "test" {
  name                  = "ddm-test"
  node_num              = 1
  time_zone             = "UTC+08:00"
  flavor_id             = var.flavor_id
  engine_id             = var.engine_id
  vpc_id                = var.vpc_id
  subnet_id             = var.subnet_id
  security_group_id     = var.sec_group_id
  enterprise_project_id = var.eps_id
  availability_zones    = [var.availability_zone]

  parameters {
    name  = "character_set_server"
    value = "utf8"
  }

  parameters {
    name  = "collation_server"
    value = "utf8_bin"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the DDM instance.
  An instance name starts with a letter, consists of `4` to `64` characters, and can contain only letters,
  digits, and hyphens (-).

* `flavor_id` - (Required, String) Specifies the ID of a product.

* `node_num` - (Required, Int) Specifies the number of nodes.

* `engine_id` - (Required, String, ForceNew) Specifies the ID of an Engine.

  Changing this parameter will create a new resource.

* `availability_zones` - (Required, List, ForceNew) Specifies the list of availability zones.

  Changing this parameter will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of a VPC.

  Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of a subnet.

  Changing this parameter will create a new resource.

* `security_group_id` - (Required, String) Specifies the ID of a security group.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project id.

  Changing this parameter will create a new resource.

* `param_group_id` - (Optional, String, ForceNew) Specifies the ID of parameter group.

  Changing this parameter will create a new resource.

* `time_zone` - (Optional, String, ForceNew) Specifies the time zone.

  Changing this parameter will create a new resource.

* `delete_rds_data` - (Optional, String) Specifies whether data stored on the associated DB instances is deleted.

* `parameters` - (Optional, List) Specifies an array of one or more parameters to be set to the instance after launched.

  The [parameters](#ddm_instance_parameters_arg) structure is documented below.

<a name="ddm_instance_parameters_arg"></a>
The `parameters` block supports:

* `name` - (Required, String) Specifies the parameter name. Some of them needs the instance to be restarted
  to take effect.
  - **bind_table**. Describes the internal data association between multiple sharded tables. It is used to inform the
    optimizer to push the join operation down to the MySQL layer for execution.
    The format is: [{tb.col1,tb2.col2},{tb.col2,tb3.col1},...].
  - **character_set_server**. Character set of the DDM server. If you need to store emoji characters,
    select **utf8mb4** and set the RDS character set to **utf8mb4** as well. When you modify the character set of the
    DDM server, the collation of the DDM server must also be changed to the corresponding value.
    + **gbk**
    + **utf8**
    + **utf8mb4**
  - **collation_server**. Character sequence of the DDM server. When modifying the character sequence of the DDM server,
    the character set of the DDM server must be changed to the corresponding value type accordingly.
    + **utf8_unicode_ci**
    + **utf8_bin**
    + **gbk_chinese_ci**
    + **gbk_bin**
    + **utf8mb4_unicode_ci**
    + **utf8mb4_bin**
  - **concurrent_execution_level**. Specifies the level of concurrent execution for sharding during logical table scanning.
    + **DATA_NODE**: Parallel scanning across different databases, with serial scanning within each shard in the same database.
    + **RDS_INSTANCE**: Parallel scanning across different RDS instances, with serial scanning within each shard in the same RDS instance.
    + **PHY_TABLE**: All physical shards are scanned in parallel.
  - **connection_idle_timeout**. The number of seconds the server waits for connection activity before closing the
    connection. The value ranges from `60` to `28800`, with a default value of `28800`, indicating that the server
    will wait for `28800` seconds before closing the connection.
  - **enable_table_recycle**. Whether to enable the table recycle bin.
    + **ON**
    + **OFF**
  - **insert_to_load_data**. Whether the insert constant value is executed using load data.
    + **ON**
    + **OFF**
  - **live_transaction_timeout_on_shutdown**. The waiting time window for in-transit transactions is measured
    in seconds, with a value range of `0-100`. The default value is `1`, indicating that the server waits for `1` second
    before closing the connection.
  - **long_query_time**. Minimum number of seconds for recording slow queries, in seconds. The value ranges
    from `0.01` to `10`. The default value is `1`, indicating that an SQL statement is defined as a slow SQL
    if its execution time is greater than or equal to `1` second.
  - **max_allowed_packet**. The maximum allowed packet size for both the server and client during a single data
    packet transmission. This value must be set as a multiple of `1024`. The range is from `1024` to `1073741824`,
    with the default value being `1073741824`.
  - **max_backend_connections**. The maximum total number of clients that each DDM node can connect to RDS
    simultaneously. The default value is `0`, which is an identifier. The actual value equals
    (Maximum number of RDS connections - 20) / Number of DDM nodes. The value range is 0 to `10000000`.
  - **max_connections**. The total number of clients that are allowed to connect simultaneously. This depends
    on the backend RDS specifications and quantity. The value ranges from `10` to `40,000`, with a default value
    of `20,000`, indicating that the total number of clients allowed to connect at the same time cannot exceed `40,000`.
  - **min_backend_connections**. The minimum total number of clients allowed to connect to RDS simultaneously on each
    DDM node. The default value is `10`. The value range is `0-10000000`.
  - **not_from_pushdown**. Whether to force pushing down query statements that do not contain FROM clauses.
    + **ON**
    + **OFF**
  - **seconds_behind_master**. The delay threshold for primary and secondary RDS nodes is measured in seconds,
    with a value range of `0-7200`. The default value is `30`, indicating that the data synchronization time between
    the primary RDS and secondary RDS should not exceed `30` seconds. If it exceeds `30` seconds, read data commands
    will not be processed by the current read node.
  - **sql_audit**. Whether enable the SQL auditing.
  - **sql_execute_timeout**. SQL execution timeout in seconds. The value ranges from `100` to `28800`, with a default
    value of `28800`, indicating that SQL execution will time out if it lasts for `28800` seconds or more.
  - **support_ddl_binlog_hint**. Add binlog hint to DDL statements.
    + **ON**
    + **OFF**
  - **ultimate_optimize**. Enable or disable the ultimate pushdown optimization feature in the optimizer.
    + **ON**
    + **OFF**
  - **transaction_policy**.
    + **XA**: XA transactions, ensuring atomicity and visibility.
    + **FREE**: Allows multiple writes, does not guarantee atomicity, and incurs no performance overhead.
    + **NO_DTX**: Single-shard transactions.

  -> **Note** `transaction_policy` is not supported by HCS **8.5.1** version.

* `value` - (Required, String) Specifies the parameter value.

-> **Note** The `parameters.name` is the key name, the `parameters.value` is the value of the key. For examlpe:

```hcl
resource "hcs_ddm_instance" "test" {
  ...

  parameters {
    name  = "character_set_server"
    value = "utf8"
  }
}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Indicates the status of the DDM instance.

* `access_ip` - Indicates the address for accessing the DDM instance.

* `access_port` - Indicates the port for accessing the DDM instance.

* `engine_version` - Indicates the engine version.

* `nodes` - Indicates the node information.
  The [nodes](#ddm_instance_nodes_attr) structure is documented below.

<a name="ddm_instance_nodes_attr"></a>
The `nodes` block supports:

* `status` - Indicates the status of the DDM instance node.

* `port` - Indicates the port of the DDM instance node.

* `ip` - Indicates the IP address of the DDM instance node.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 15 minutes.

* `update` - Default is 60 minutes.

* `delete` - Default is 10 minutes.

## Import

The ddm instance can be imported using the `id`, e.g.

```
$ terraform import hcs_ddm_instance.test 4bc36477c36642479acf2d90751c8c29in09
```
