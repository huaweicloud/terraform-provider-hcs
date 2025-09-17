---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_dms_rocketmq_instance"
description: |-
  Manage DMS RocketMQ instance resources within HuaweiCloudStack.
---

# hcs_dms_rocketmq_instance

Manage DMS RocketMQ instance resources within HuaweiCloudStack.

## Example Usage

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}
variable "availability_zones" {
  type = list(string)
}

resource "hcs_dms_rocketmq_instance" "test" {
  name               = "rocketmq_name_test"
  description        = "this is a rocketmq instance"
  engine_version     = "4.8.0"
  storage_space      = 300
  vpc_id             = var.vpc_id
  subnet_id          = var.subnet_id
  security_group_id  = var.security_group_id
  availability_zones = var.availability_zones
  flavor_id          = "c6.4u8g.cluster"
  storage_spec_code  = "dms.physical.storage.high.v2"
  broker_num         = 1
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the DMS RocketMQ instance.
  An instance name starts with a letter, consists of 4 to 64 characters, and can contain only letters,
  digits, underscores (_), and hyphens (-).

* `engine_version` - (Required, String, ForceNew) Specifies the version of the RocketMQ engine. Value: **5.x**.

  Changing this parameter will create a new resource.

* `storage_space` - (Required, Int, ForceNew) Specifies the message storage capacity, Unit: GB.
  Value range: 300-3000.

  Changing this parameter will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of a VPC.

  Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of a subnet.

  Changing this parameter will create a new resource.

* `security_group_id` - (Required, String) Specifies the ID of a security group.

* `availability_zones` - (Required, List, ForceNew) Specifies the list of availability zone names, where
  instance brokers reside and which has available resources.

  Changing this parameter will create a new resource.

* `flavor_id` - (Required, String, ForceNew) Specifies a product ID. The options are as follows:
  + **c6.2u8g.single.x86** or **c6.2u8g.single.arm**. maximum number of topics on each broker: 50;
    maximum number of consumer groups on each broker: 100.
  + **c6.4u16g.single.x86** or **c6.4u16g.single.arm**. maximum number of topics on each broker: 100;
    maximum number of consumer groups on each broker: 200.
  + **c6.8u32g.single.x86** or **c6.8u32g.single.arm**. maximum number of topics on each broker: 200;
    maximum number of consumer groups on each broker: 400.
  + **c6.16u64g.single.x86** or **c6.16u64g.single.arm**. maximum number of topics on each broker: 300;
    maximum number of consumer groups on each broker: 600.
  + **c6.32u128g.single.x86** or **c6.32u128g.single.arm**. maximum number of topics on each broker: 400;
    maximum number of consumer groups on each broker: 800.

  Changing this parameter will create a new resource.

* `storage_spec_code` - (Required, String, ForceNew) Specifies the storage I/O specification.
  The options are as follows:
  + **dms.physical.storage.high.v2**: high I/O disk
  + **dms.physical.storage.ultra.v2**: ultra-high I/O disk

  Changing this parameter will create a new resource.

* `broker_num` - (Required, Int, ForceNew) Specifies the broker numbers.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of the DMS RocketMQ instance.

  The description can contain a maximum of 1024 characters.

* `ssl_enable` - (Optional, Bool, ForceNew) Specifies whether the RocketMQ **SASL_SSL** is enabled.
  Defaults to **false**.

  Changing this parameter will create a new resource.

* `ipv6_enable` - (Optional, Bool, ForceNew) Specifies whether to support IPv6. Defaults to **false**.

  Changing this parameter will create a new resource.

* `enable_publicip` - (Optional, Bool, ForceNew) Specifies whether to enable public access. Defaults to **false**.

  Changing this parameter will create a new resource.

* `publicip_id` - (Optional, String, ForceNew) Specifies the ID of the EIP bound to the instance.
  Use commas (,) to separate multiple EIP IDs. This parameter is **Required** when `enable_publicip` set to **true**.

  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project id of the instance.

  Changing this parameter will create a new resource.

* `enable_acl` - (Optional, Bool) Specifies whether access control is enabled. Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `status` - Indicates the status of the DMS RocketMQ instance.

* `type` - Indicates the DMS RocketMQ instance type. Value: cluster.

* `specification` - Indicates the instance specification. For a cluster DMS RocketMQ instance, VM specifications
  and the number of nodes are returned.

* `maintain_begin` - Indicates the time at which the maintenance window starts. The format is HH:mm:ss.

* `maintain_end` - Indicates the time at which the maintenance window ends. The format is HH:mm:ss.

* `used_storage_space` - Indicates the used message storage space. Unit: GB.

* `publicip_address` - Indicates the public IP address.

* `cross_vpc_info` - Indicates the Cross-VPC access information.

* `node_num` - Indicates the node quantity.

* `new_spec_billing_enable` - Indicates whether billing based on new specifications is enabled.

* `namesrv_address` - Indicates the metadata address.

* `broker_address` - Indicates the service data address.

* `public_namesrv_address` - Indicates the public network metadata address.

* `public_broker_address` - Indicates the public network service data address.

* `resource_spec_code` - Indicates the resource specifications.

* `charging_mode` - Indicates the charging mode of the instance.

* `vpc_name` - Indicates the name of VPC.

* `subnet_name` - Indicates the name of subnet.

* `security_group_name` - Indicates the name of security group.

* `user_id` - Indicates the user ID who created the instance.

* `user_name` - Indicates the user-name who created the instance.

* `created_at` - Indicates when the instance is created.

* `enable_log_collection` - Whether to enable log collection function.

* `storage_resource_id` - Indicates ID of storage resource.

* `service_type` - Indicates the type of service. Such as **advanced**.

* `storage_type` - Indicates the type of storage. Such as **hec**.

* `extend_times` - Indicates the time when the instance extended.

* `support_features` - Indicates features supported by the instance.

* `disk_encrypted` - Whether to enable disk encryption.

* `ces_version` - Indicates the version of CES.

* `grpc_address` - Indicates the gRPC connection address (this field is only displayed for RocketMQ version 5.x).

* `public_grpc_address` - Indicates Public network gRPC connection address (this field is only displayed for
  RocketMQ version 5.x).

* `cross_vpc_accesses` - Indicates the Access information of cross-VPC.
  The [cross_vpc_accesses](#rocketmq_instance) structure is documented below.

<a name="rocketmq_instance"></a>
The `cross_vpc_accesses` block supports:

* `advertised_ip` - The advertised IP Address or domain name.

* `listener_ip` - The listener IP address.

* `port` - The port number.

* `port_id` - The port ID associated with the address.

## Import

The rocketmq instance can be imported using the `id`, e.g.

```
$ terraform import hcs_dms_rocketmq_instance.test 8d3c7938-dc47-4937-a30f-c80de381c5e3
```
