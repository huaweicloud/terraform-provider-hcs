---
subcategory: "ROMA Connect"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_roma_connect_instance"
description: ""
---

# hcs_roma_connect_instance

Manage a ROMA Connect instance resource within HuaweiCloudStack.

## Example Usage

### Create a Family Bucket instance using kafka engine 

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}

resource "hcs_roma_connect_instance" "test1" {
  name              = "roma-test"
  description       = "terraform test description"
  product_id        = "00300-30101-0--0"
  available_zones   = ["az0.dc0"]
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.security_group_id
  ipv6_enable       = false
  enable_all        = true
  cpu_architecture  = "arm"
  maintain_begin    = "22:00:00"
  maintain_end      = "02:00:00"

  mqs {
    connector_enable = false
    enable_publicip  = false
    engine_version   = "2.7"
    retention_policy = "produce_reject"
    ssl_enable       = true
    vpc_client_plain = false
    trace_enable     = false
  }
}
```

### Create a MQS instance using RocketMQ engine

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "security_group_id" {}

resource "hcs_roma_connect_instance" "test1" {
  name              = "roma-test"
  description       = "terraform test description"
  product_id        = "00300-30101-0--0"
  available_zones   = ["az0.dc0"]
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.security_group_id
  ipv6_enable       = false
  enable_all        = false
  cpu_architecture  = "arm"
  maintain_begin    = "22:00:00"
  maintain_end      = "02:00:00"

  mqs {
    rocketmq_enable = true
    ssl_enable      = true
    enable_acl      = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the instance. If omitted,
  the provider-level region will be used. Changing this will create a new instance.

* `name` - (Required, String, ForceNew) Specifies the name of instance.
  Changing this will create a new instance.

* `description` - (Required, String, ForceNew) Specifies the description of instance.
  Changing this will create a new instance.

* `product_id` - (Required, String, ForceNew) Specifies the product id of instance.
  - **00300-30101-0--0**. Create a **basic** instance.
  - **00300-30102-0--0**. Create a **professional** instance.
  - **00300-30103-0--0**. Create a **enterprise** instance.
  - **00300-30104-0--0**. Create a **platinum** instance.
  - **00300-30105-0--0**. Create a **platinumXB** instance.

  Changing this will create a new instance.

* `available_zones` - (Required, List, ForceNew) Specifies the availability zone name.
  Changing this will create a new instance.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of the vpc used to create the instance.
  Changing this will create a new instance.

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of the subnet used to create the instance.
  Changing this will create a new instance.

* `security_group_id` - (Required, String, ForceNew) Specifies the ID of the security group used to create the instance.
  Changing this will create a new instance.

* `ipv6_enable` - (Required, Bool, ForceNew) Specifies whether using IPv6 to create instance. If true, the subnet must
  supports IPV6. Changing this will create a new instance.

* `enable_all` - (Required, Bool, ForceNew) Specifies whether a bucket instance is a family bucket instance.
  - **true**. Create a Family Bucket instance.
  - **false**. Create a MQS instance.

  Changing this will create a new instance.

* `cpu_architecture` - (Required, String, ForceNew) Architecture type, the valid value are as follows.
  - **x86**. Create a x86 architecture instance.
  - **arm**. Create a arm architecture instance.

  Changing this will create a new instance.

* `mqs` - (Required, Map, ForceNew) Specifies MQS service parameters. The [mqs](#roma_mqs) structure is documented
  below. Changing this will create a new instance.

* `maintain_begin` - (Required, String, ForceNew) Maintain start time. Example **22:00:00**.
  Changing this will create a new instance.

* `maintain_end` - (Required, String, ForceNew) Maintain end time. Example **02:00:00**.
  Changing this will create a new instance.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID.
  Changing this will create a new instance.

* `entrance_bandwidth_size` - (Optional, Int, ForceNew) Specifies the ingress bandwidth. This parameter is **Required**
  when `enable_publicip` is true.

  Changing this will create a new instance.

<a name="roma_mqs"></a>
The `mqs` object supports the following:

* `enable_publicip` - (Optional, Bool, ForceNew) Specifies whether public network access is supported.
  Default to **false**. Changing this will create a new instance.

* `ssl_enable` - (Optional, Bool, ForceNew) Specifies whether public network access is supported.
  Default to **false**. Changing this will create a new instance.

* `engine_version` - (Optional, String, ForceNew) Specifies the **Kafka** engine version. This parameter is conflicted
  with `rocketmq_enable`. Valid version are **2.7**, **2.3.0**, **1.1.0**.

  Changing this will create a new instance.

* `rocketmq_enable` - (Optional, Bool, ForceNew) Specifies create a **RocketMQ** instance. This parameter is conflicted
  with `engine_version`. Default to **false**.

  Changing this will create a new instance.

* `enable_acl` - (Optional, Bool, ForceNew) Specifies whether to enable acl. This parameter is valid only when the
  `rocketmq_enable` is true.

  Changing this will create a new instance.

* `retention_policy` - (Optional, String, ForceNew) Message integration capacity threshold policy. This parameter is
  valid only when the engine type is **Kafka**.
  - **produce_reject**: production is restricted.
  - **time_base**: automatic deletion.

  Changing this will create a new instance.

* `trace_enable` - (Optional, Bool, ForceNew) Specifies whether to enable message tracing.
  This parameter is valid only when the engine type is **Kafka**.

  Changing this will create a new instance.

* `vpc_client_plain` - (Optional, Bool, ForceNew) Specifies whether to enable plaintext access in a VPC. This
  parameter is valid only when the engine type is **Kafka**. Default to **false**.

  Changing this will create a new instance.

* `connector_enable` - (Optional, Bool, ForceNew) Specifies whether to enable smart connect. This parameter is valid
  only when the engine type is **Kafka**.

  Changing this will create a new instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ROMA Connect instance ID.

* `site_id` - ID of a site instance. This parameter is valid only in the site scenario.

* `flavor_id` - Specifies the instance flavor ID.

* `flavor_type` - Specifies the instance flavor type. The type are **basic**, **professional**, **enterprise**,
  **platinum**, **platinumXB**.

* `cpu_arch` - CPU architecture type. The options are as follows.
  - **x86_64**: x86 architecture. 
  - **aarch64**: ARM architecture.

* `publicip_id` - Specifies the EIP ID bound to the instance.

* `publicip_address` - Specifies the EIP bound to the instance.

* `publicip_enable` - Indicates whether to enable public network access. When public network access is enabled,
  **publicip_id** is mandatory.

* `connect_address` - ROMA Connect connection address.

* `charge_type` - Charging type of a resource specification. postPaid: indicates the pay-per-use charging type.
  The service is used first and then paid.

* `bandwidths` - Specifies the egress bandwidth of the public network.

* `ipv6_connect_address` - ROMA Connect connection address in IPv6 scenarios.

* `external_elb_enable` - Indicates whether to enable the external ELB.

* `external_elb_id` - IP address ID of the external ELB.

* `external_elb_address` - IP address of the external ELB.

* `external_eip_bound` - Specifies whether the external ELB is bound to a public IP address.

* `external_eip_id` - Specifies the ID of the public network address bound to the external load balancer.

* `external_eip_address` - Specifies the public IP address bound to the external ELB.

* `create_time` - Creation time.

* `update_time` - Last update time.

* `resources` - The description of resource information. The [resources](#attr_resources) structure is documented below.

* <a name="attr_resources"></a>
The `resources` object supports the following:

* `mqs` - Specifies MQS service parameters. The [mqs](#attr_mqs) structure is documented below.

<a name="attr_mqs"></a>
The `mqs` object supports the following:

* `id` - Indicates the instance ID.

* `enable` - Indicates whether to enable the MQS component.

* `retention_policy` - Message integration capacity threshold policy.

* `ssl_enable` - Whether to enable SSL.

* `trace_enable` - Whether to enable the message track.

* `vpc_client_plain` - Whether to enable plaintext access in a VPC.

* `partition_num` - Number of partitions.

* `specification` - Indicates the specification.

* `private_connect_address` - Internal connection address.

* `public_connect_address` - External connection address.

* `private_restful_address` - Internal restful connection address.

* `public_restful_address` - Internal restful connection address.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

* `delete` - Default is 30 minutes.

## Import

The instance can be imported using the instance ID.

```bash
$ terraform import hcs_roma_connect_instance.test <instance_id>
```
