---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_dcs_instance"
description: |-
  Manages a DCS instance within HuaweiCloudStack.
---

# hcs_dcs_instance

Manages a DCS instance within HuaweiCloudStack.

~> **WARNING:** DCS for Memcached is about to become unavailable and is no longer sold in some regions.
You can use DCS for Redis **3.0**, **4.0**, **5.0** or **6.0** instead.
It is not possible to create Memcached instances through this resource.

## Example Usage

### Create a single mode Redis instance

```hcl
variable vpc_id {}
variable subnet_id {}
variable pwd {}
variable az {}

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
  availability_zones = [var.az]
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
variable az1 {}
variable az2 {}

resource "hcs_dcs_instance" "instance_2" {
  name               = "redis_name"
  engine             = "Redis"
  engine_version     = "5.0"
  capacity           = 4
  flavor             = "redis.ha.xu1.large.r2.4"
  # The first is the primary availability zone (var.az1),
  # and the second is the standby availability zone (var.az2).
  availability_zones = [var.az1, var.az2]
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
  If omitted, the provider-level region will be used. Changing this parameter will create a new instance.

* `name` - (Required, String) Specifies the name of an instance.
  The name must be 4 to 64 characters and start with a letter.
  Only chinese, letters (case-insensitive), digits, underscores (_) ,and hyphens (-) are allowed.

* `engine` - (Required, String, ForceNew) Specifies a cache engine. The valid value is **Redis**.

  Changing this parameter will create a new instance.

* `engine_version` - (Optional, String, ForceNew) Specifies the version of a cache engine.
  It is mandatory when the engine is **Redis**, the value can be **3.0**, **4.0**, **5.0** or **6.0**.

  Changing this parameter will create a new instance.

* `capacity` - (Required, Float) Specifies the cache capacity. Unit: GB.
  + **Redis3.0**
    - **Single-Node** and **Master/Standby** type instance values: `2`, `4`, `8`, `16`, `32` and `64`.
    - **Proxy cluster** instance specifications values `64`, `128`, `256`, `512`, and `1024`.
  + **Redis4.0, Redis5.0 and Redis6.0**
    - **Single-Node** and **Master/Standby** type instance values: `0.125`, `0.25`, `0.5`, `1`, `2`, `4`, `8`, `16`,
      `32` and `64`.
    - **Cluster** instance values `4`,`8`,`16`, `24`, `32`, `48`, `64`, `96`, `128`, `192`, `256`, `384`, `512`,
      `768` and `1024`.
    - **Read/Write Separation** instance values `8`, `16`, `32` and `64`.

* `flavor` - (Required, String) Specifies the flavor of the Redis instance.  
  You can query the **flavor** as follows:
  + It can be obtained through this data source `hcs_dcs_flavors`.
  + Log in to the DCS console, click **Buy DCS Instance**, and find the corresponding instance specification.

* `availability_zones` - (Required, List, ForceNew) Specifies the code of the AZ where the cache node resides.
  **Master/Standby**, **Proxy Cluster**, and **Redis Cluster** instances support cross-AZ deployment.
  You can specify an AZ for the standby node. When specifying AZs for nodes, use commas (,) to separate AZs.

  Changing this parameter will create a new instance.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of VPC which the instance belongs to.

  Changing this parameter will create a new instance.

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of subnet which the instance belongs to.

  Changing this parameter will create a new instance.

* `security_group_id` - (Optional, String) Specifies the ID of the security group which the instance belongs to.
  This parameter is **Required** for Redis **3.0** version.

* `private_ip` - (Optional, String, ForceNew) Specifies the IP address of the instance,
  which can only be the currently available IP address the selected subnet.
  You can specify an available IP for the Redis instance (except for the **Redis Cluster** type).
  If omitted, the system will automatically allocate an available IP address to the Redis instance.

  Changing this parameter will create a new instance.

* `template_id` - (Optional, String, ForceNew) Specifies the Parameter Template ID.

  Changing this parameter will create a new instance.

* `port` - (Optional, Int) Specifies the port customization. Defaults to **6379**.

  -> **Note** This parameter is only supported by Redis **4.0** and **later** version.

* `password` - (Optional, String) Specifies the password of a DCS instance.
  The password of a DCS instance must meet the following complexity requirements:
  + Must be a string of `8 to 32` bits in length.
  + Must contain three combinations of the following four characters: Lower case letters, uppercase letter, digital,
    Special characters include (`~!@#$^&*()-_=+\\|{}:,<.>/?).
  + The new password cannot be the same as the old password.

* `whitelists` - (Optional, List) Specifies the IP addresses which can access the instance.
  The [whitelists](#dcs_whitelists) object structure is documented below.

  -> **Note** This parameter is only supported by Redis **4.0** and **later** version.

* `whitelist_enable` - (Optional, Bool) Whether enable or disable the IP address whitelists. Defaults to **true**.
  If the whitelist is disabled, all IP addresses connected to the VPC can access the instance.

* `maintain_begin` - (Optional, String) Specifies time at which the maintenance time window starts.
  + The start time and end time of a maintenance time window must indicate the time segment of a supported maintenance
    time window.
  + The start time must be on the hour, such as **18:00:00**.
  + If parameter `maintain_begin` is left blank, parameter `maintain_end` is also blank.

* `maintain_end` - (Optional, String) Specifies time at which the maintenance time window ends.
  + The start time and end time of a maintenance time window must indicate the time segment of a supported maintenance
    time window.
  + The end time is one hour later than the start time. For example, if the start time is **18:00:00**, the end time is
    **22:00:00**.
  + If parameter `maintain_end` is left blank, parameter `maintain_begin` is also blank.

-> **NOTE:** Parameters `maintain_begin` and `maintain_end` must be set in pairs.

* `backup_policy` - (Optional, List) Specifies the backup configuration to be used with the instance.
  The [backup_policy](#dcs_backup_policy) object structure is documented below.

  -> **NOTE:** This parameter is not supported when the instance type is single.

* `parameters` - (Optional, List) Specifies an array of one or more parameters to be set to the DCS instance after
  launched. You can check on console to see which parameters supported.
  The [parameters](#dcs_parameters) object structure is documented below.

* `rename_commands` - (Optional, Map) Specifies the critical command renaming. The valid values are as follows:
  + **command**
  + **keys**
  + **flushdb**
  + **flushall**
  + **scan**
  + **hscan**
  + **sscan**
  + **zscan**
  + **hgetall**

  -> **Note** This parameter is only supported by Redis **4.0** and **later** version.

* `tags` - (Optional, Map) The key/value pairs to associate with the dcs instance.

* `access_user` - (Optional, String, ForceNew) Specifies the username. If the cache engine is *Redis*, you do not need
  to set this parameter. The username starts with a letter, consists of 1 to 64 characters, and supports only letters,
  digits, and hyphens (-).

  Changing this parameter will create a new instance.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project id of the dcs instance.

  Changing this parameter will create a new instance.

* `description` - (Optional, String) Specifies the description of an instance.
  It is a string that contains a maximum of **1024** characters.

* `deleted_nodes` - (Optional, List) Specifies the ID of the replica to delete. 
  Currently, only one replica can be deleted at a time.

  -> **Note** This parameter is mandatory when you delete replicas of a master/standby DCS redis **4.0** and
  **later** version.

* `reserved_ips` - (Optional, List) Specifies IP addresses to retain. Mandatory during cluster scale-in. If this
  parameter is not set, the system randomly deletes unnecessary shards.

<a name="dcs_whitelists"></a>
The `whitelists` block supports:

* `group_name` - (Required, String) Specifies the name of IP address group.

* `ip_address` - (Required, List) Specifies the list of IP address or CIDR which can be whitelisted for an instance.
  The maximum is **20**.

<a name="dcs_backup_policy"></a>
The `backup_policy` block supports:

* `backup_type` - (Optional, String) Specifies the backup type. Defaults to `auto`. The valid values are as follows:
  + `auto`: automatic backup.
  + `manual`: manual backup.

* `save_days` - (Optional, Int) Specifies the retention time. Unit: day, the value ranges `from 1 to 7`.
  This parameter is **required** if the `backup_type` is **auto**.

* `period_type` - (Optional, String) Specifies the interval at which backup is performed. Defaults to `weekly`.
  Currently, only weekly backup is supported.

* `backup_at` - (Required, List) Specifies the day in a week on which backup starts, the value ranges `from 1 to 7`.
  + **1**: Monday
  + **2**: Tuesday
  + **3**: Wednesday
  + **4**: Thursday
  + **5**: Friday
  + **6**: Saturday
  + **7**: Sunday

* `begin_at` - (Required, String) Specifies the time at which backup starts.
  **00:00-01:00** indicates that the backup starts at midnight. It can only be set for whole-hour time periods,
  with a minimum interval of one hour.

<a name="dcs_parameters"></a>
The `parameters` block supports:

* `id` - (Required, String) Specifies the ID of the configuration item.

* `name` - (Required, String) Specifies the name of the configuration item.

* `value` - (Required, String) Specifies the value of the configuration item.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `status` - The cache instance status. The valid values are as follows:
  + `RUNNING`: The instance is running properly.
    Only instances in the Running state can provide in-memory cache service.
  + `ERROR`: The instance is not running properly.
  + `RESTARTING`: The instance is being restarted.
  + `EXTENDING`: The instance is being scaled up.
  + `RESTORING`: The instance data is being restored.
  + `FLUSHING`: The DCS instance is being cleared.

* `domain_name` - The domain name of the instance. Usually, we use domain name and port to connect to the DCS instances.

* `max_memory` - The total memory size. Unit: MB.

* `used_memory` - The size of the used memory. Unit: MB.

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
