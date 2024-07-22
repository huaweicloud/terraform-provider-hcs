---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_dcs_instance"
description: ""
---

# hcs_dcs_instance

Manages a DCS instance within HuaweiCloudStack.

!> **WARNING:** DCS for Memcached is about to become unavailable and is no longer sold in some regions.
You can use DCS for Redis 4.0, 5.0 or 6.0 instead. It is not possible to create Memcached instances through this resource.

## Example Usage

### Create a single mode Redis instance

```hcl
variable vpc_id {}
variable subnet_id {}
variable pwd {}

data "hcs_dcs_flavors" "single_flavors" {
  cache_mode = "single"
  capacity   = 0.125
}

resource "hcs_dcs_instance" "instance_1" {
  name               = "redis_single_instance"
  engine             = "Redis"
  engine_version     = "5.0"
  capacity           = data.hcs_dcs_flavors.single_flavors.capacity
  flavor             = data.hcs_dcs_flavors.single_flavors.flavors[0].name
  availability_zones = ["cn-north-1a"]
  password           = var.pwd
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id
}
```

### Create Master/Standby mode Redis instances with backup policy

```hcl
variable vpc_id {}
variable subnet_id {}
variable pwd {}

resource "hcs_dcs_instance" "instance_2" {
  name               = "redis_name"
  engine             = "Redis"
  engine_version     = "5.0"
  capacity           = 4
  flavor             = "redis.ha.xu1.large.r2.4"
  # The first is the primary availability zone (cn-north-1a),
  # and the second is the standby availability zone (cn-north-1b).
  availability_zones = ["cn-north-1a", "cn-north-1b"]
  password           = var.pwd
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id

  backup_policy {
    backup_type = "auto"
    save_days   = 3
    backup_at   = [1, 3, 5, 7]
    begin_at    = "02:00-04:00"
  }

  whitelists {
    group_name = "test-group1"
    ip_address = ["192.168.10.100", "192.168.0.0/24"]
  }
  whitelists {
    group_name = "test-group2"
    ip_address = ["172.16.10.100", "172.16.0.0/24"]
  }

  parameters {
    id    = "1"
    name  = "timeout"
    value = "500"
  }
  parameters {
    id    = "3"
    name  = "hash-max-ziplist-entries"
    value = "4096"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the DCS instance resource.
  If omitted, the provider-level region will be used. Changing this creates a new DCS instance resource.

* `name` - (Required, String) Specifies the name of an instance.
  The name must be 4 to 64 characters and start with a letter.
  Only chinese, letters (case-insensitive), digits, underscores (_) ,and hyphens (-) are allowed.

* `engine` - (Required, String, ForceNew) Specifies a cache engine. Options: *Redis*.
  Changing this creates a new instance.

* `engine_version` - (Optional, String, ForceNew) Specifies the version of a cache engine.
  It is mandatory when the engine is *Redis*, the value can be 3.0, 4.0, 5.0 or 6.0.
  Changing this creates a new instance.

* `capacity` - (Required, Float) Specifies the cache capacity. Unit: GB.
  + **Redis4.0, Redis5.0 and Redis6.0**: Stand-alone and active/standby type instance values: `0.125`, `0.25`,
    `0.5`, `1`, `2`, `4`, `8`, `16`, `32` and `64`.
    Cluster instance specifications support `4`,`8`,`16`, `24`, `32`, `48`, `64`, `96`, `128`, `192`, `256`,
    `384`, `512`, `768` and `1024`.
  + **Redis3.0**: Stand-alone and active/standby type instance values: `2`, `4`, `8`, `16`, `32` and `64`.
    Proxy cluster instance specifications support `64`, `128`, `256`, `512`, and `1024`.

* `flavor` - (Required, String) The flavor of the cache instance, which including the total memory, available memory,
  maximum number of connections allowed, maximum/assured bandwidth and reference performance.
  It also includes the modes of Redis instances. You can query the *flavor* as follows:
  + It can be obtained through this data source `hcs_dcs_flavors`.
  + Log in to the DCS console, click *Buy DCS Instance*, and find the corresponding instance specification.

* `availability_zones` - (Required, List, ForceNew) The code of the AZ where the cache node resides.
  Master/Standby, Proxy Cluster, and Redis Cluster DCS instances support cross-AZ deployment.
  You can specify an AZ for the standby node. When specifying AZs for nodes, use commas (,) to separate AZs.
  Changing this creates a new instance.

* `vpc_id` - (Required, String, ForceNew) The ID of VPC which the instance belongs to.
  Changing this creates a new instance resource.

* `subnet_id` - (Required, String, ForceNew) The ID of subnet which the instance belongs to.
  Changing this creates a new instance resource.

* `security_group_id` - (Optional, String) The ID of the security group which the instance belongs to.
  This parameter is mandatory for Redis 3.0 version.

* `private_ip` - (Optional, String, ForceNew) The IP address of the DCS instance,
  which can only be the currently available IP address the selected subnet.
  You can specify an available IP for the Redis instance (except for the Redis Cluster type).
  If omitted, the system will automatically allocate an available IP address to the Redis instance.
  Changing this creates a new instance resource.

* `template_id` - (Optional, String, ForceNew) The Parameter Template ID.
  Changing this creates a new instance resource.

* `port` - (Optional, Int) Port customization, which is supported only by Redis 4.0 and Redis 5.0 instances.
  Redis instance defaults to 6379.

* `password` - (Optional, String) Specifies the password of a DCS instance.
  The password of a DCS instance must meet the following complexity requirements:
  + Must be a string of 8 to 32 bits in length.
  + Must contain three combinations of the following four characters: Lower case letters, uppercase letter, digital,
    Special characters include (`~!@#$^&*()-_=+\\|{}:,<.>/?).
  + The new password cannot be the same as the old password.

* `whitelists` - (Optional, List) Specifies the IP addresses which can access the instance.
  This parameter is valid for Redis 4.0 and 5.0 versions. The structure is described below.

* `whitelist_enable` - (Optional, Bool) Enable or disable the IP address whitelists. Defaults to true.
  If the whitelist is disabled, all IP addresses connected to the VPC can access the instance.

* `maintain_begin` - (Optional, String) Time at which the maintenance time window starts. Defaults to **02:00:00**.
  + The start time and end time of a maintenance time window must indicate the time segment of a supported maintenance
    time window.
  + The start time must be on the hour, such as **18:00:00**.
  + If parameter `maintain_begin` is left blank, parameter `maintain_end` is also blank.
    In this case, the system automatically allocates the default start time **02:00:00**.

* `maintain_end` - (Optional, String) Time at which the maintenance time window ends. Defaults to **06:00:00**.
  + The start time and end time of a maintenance time window must indicate the time segment of a supported maintenance
    time window.
  + The end time is one hour later than the start time. For example, if the start time is **18:00:00**, the end time is
    **22:00:00**.
  + If parameter `maintain_end` is left blank, parameter `maintain_begin` is also blank.
    In this case, the system automatically allocates the default end time **06:00:00**.

-> **NOTE:** Parameters `maintain_begin` and `maintain_end` must be set in pairs.

* `backup_policy` - (Optional, List) Specifies the backup configuration to be used with the instance.
  The structure is described below.

  -> **NOTE:** This parameter is not supported when the instance type is single.

* `parameters` - (Optional, List) Specify an array of one or more parameters to be set to the DCS instance after
  launched. You can check on console to see which parameters supported.
  The [parameters](#DcsInstance_Parameters) structure is documented below.

* `rename_commands` - (Optional, Map) Critical command renaming, which is supported only by Redis 4.0 and
  Redis 5.0 instances but not by Redis 3.0 instance.
  The valid commands that can be renamed are: **command**, **keys**, **flushdb**, **flushall** and **hgetall**.

* `tags` - (Optional, Map) The key/value pairs to associate with the dcs instance.

* `access_user` - (Optional, String, ForceNew) Specifies the username. If the cache engine is *Redis*, you do not need
  to set this parameter. The username starts with a letter, consists of 1 to 64 characters, and supports only letters,
  digits, and hyphens (-). Changing this creates a new instance.

* `description` - (Optional, String) Specifies the description of an instance.
  It is a string that contains a maximum of 1024 characters.

* `deleted_nodes` - (Optional, List) Specifies the ID of the replica to delete. This parameter is mandatory when
  you delete replicas of a master/standby DCS Redis 4.0 or 5.0 instance. Currently, only one replica can be deleted
  at a time.

* `reserved_ips` - (Optional, List) Specifies IP addresses to retain. Mandatory during cluster scale-in. If this
  parameter is not set, the system randomly deletes unnecessary shards.

The `whitelists` block supports:

* `group_name` - (Required, String) Specifies the name of IP address group.

* `ip_address` - (Required, List) Specifies the list of IP address or CIDR which can be whitelisted for an instance.
  The maximum is 20.

The `backup_policy` block supports:

* `backup_type` - (Optional, String) Backup type. Default value is `auto`. The valid values are as follows:
  + `auto`: automatic backup.
  + `manual`: manual backup.

* `save_days` - (Optional, Int) Retention time. Unit: day, the value ranges from 1 to 7.
  This parameter is required if the backup_type is **auto**.

* `period_type` - (Optional, String) Interval at which backup is performed. Default value is `weekly`.
  Currently, only weekly backup is supported.

* `backup_at` - (Required, List) Day in a week on which backup starts, the value ranges from 1 to 7.
  Where: 1 indicates Monday; 7 indicates Sunday.

* `begin_at` - (Required, String) Time at which backup starts.
  Format: `hh24:00-hh24:00`, "00:00-01:00" indicates that backup starts at 00:00:00.

<a name="DcsInstance_Parameters"></a>
The `parameters` block supports:

* `id` - (Required, String) Specifies the ID of the configuration item.

* `name` - (Required, String) Specifies the name of the configuration item.

* `value` - (Required, String) Specifies the value of the configuration item.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID in UUID format.

* `status` - Cache instance status. The valid values are as follows:
  + `RUNNING`: The instance is running properly.
    Only instances in the Running state can provide in-memory cache service.
  + `ERROR`: The instance is not running properly.
  + `RESTARTING`: The instance is being restarted.
  + `EXTENDING`: The instance is being scaled up.
  + `RESTORING`: The instance data is being restored.
  + `FLUSHING`: The DCS instance is being cleared.

* `domain_name` - Domain name of the instance. Usually, we use domain name and port to connect to the DCS instances.

* `max_memory` - Total memory size. Unit: MB.

* `used_memory` - Size of the used memory. Unit: MB.

* `vpc_name` - The name of VPC which the instance belongs to.

* `subnet_name` - The name of subnet which the instance belongs to.


* `security_group_name` - The name of security group which the instance belongs to.

* `order_id` - The ID of the order that created the instance.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 120 minutes.
* `update` - Default is 120 minutes.
* `delete` - Default is 15 minutes.

## Import

DCS instance can be imported using the `id`, e.g.

```bash
terraform import hcs_dcs_instance.instance_1 80e373f9-872e-4046-aae9-ccd9ddc55511
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `password`, `auto_renew`, `period`, `period_unit`, `rename_commands`,
`internal_version`, `save_days`, `backup_type`, `begin_at`, `period_type`, `backup_at`, `parameters`.
It is generally recommended running `terraform plan` after importing an instance.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```
resource "hcs_dcs_instance" "instance_1" {
    ...

  lifecycle {
    ignore_changes = [
      password, rename_commands,
    ]
  }
}
```
