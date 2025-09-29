---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_dms_rocketmq_topic"
description: |-
  Manages DMS RocketMQ topic resources within HuaweiCloudStack.
---

# hcs_dms_rocketmq_topic

Manages DMS RocketMQ topic resources within HuaweiCloudStack.

## Example Usage

```hcl
variable "instance_id" {}

resource "hcs_dms_rocketmq_topic" "test" {
  instance_id = var.instance_id
  name        = "topic_test"
  permission  = "all"

  queue_num = 8

  queues  {
    broker    = "broker-0"
    queue_num = 3
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the rocketMQ instance.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the topic.

  Changing this parameter will create a new resource.

* `queue_num` - (Optional, Int, ForceNew) Specifies the number of queues. Default to 8.

  Changing this parameter will create a new resource.

* `permission` - (Optional, String, ForceNew) Specifies the permissions of the topic.
  Value options: **all**, **sub**, **pub**. Default to **all**.

  Changing this parameter will create a new resource.

* `queues` - (Optional, List, ForceNew) Specifies the queues information of the topic.
  It's only valid when RocketMQ instance version is **4.8.0**.  
  The [queues](#rocketmq_tpoic_queues) structure is documented below.
  
  Changing this parameter will create a new resource.

  -> This parameter supported only in HCS **8.5.1** and **later** version.

<a name="rocketmq_tpoic_queues"></a>
The `queues` block supports:

* `broker` - (Optional, String) Specifies the associated broker.

* `queue_num` - (Optional, Int) Specifies the number of the queues.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `total_read_queue_num` - The total number of read queues.

* `total_write_queue_num` - The total number of write queues.

* `brokers` - The list of associated brokers of the topic.
  The [brokers](#rocketmq_topic_brokers) structure is documented below.

  -> It's only valid when RocketMQ instance version is **4.8.0**.

<a name="rocketmq_topic_brokers"></a>
The `brokers` block supports:

* `name` - Indicates the name of the broker.

* `read_queue_num` - Indicates the read queues number of the broker.

* `write_queue_num` - Indicates the read queues number of the broker.

## Import

The rocketmq topic can be imported using the rocketMQ instance ID and topic name separated by a slash, e.g.

```
$ terraform import hcs_dms_rocketmq_topic.test c8057fe5-23a8-46ef-ad83-c0055b4e0c5c/topic_1
```
