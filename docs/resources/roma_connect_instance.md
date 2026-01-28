---
subcategory: "ROMA Connect"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_roma_connect_instance"
description: |-
  Manage a ROMA Connect instance resource within HuaweiCloudStack.
---

# hcs_roma_connect_instance

Manage a ROMA Connect instance resource within HuaweiCloudStack.

~> **WARNING:** To use this resource, you need to manually register the ROMA Connect API in your environment.

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

resource "hcs_roma_connect_instance" "test2" {
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
  the provider-level region will be used.

  Changing this will create a new instance.

* `name` - (Required, String, ForceNew) Specifies the name of instance.

  Changing this will create a new instance.

* `description` - (Optional, String, ForceNew) Specifies the description of instance.

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

* `ipv6_enable` - (Required, Bool, ForceNew) Specifies whether using IPv6 to create instance. 

  Changing this will create a new instance.

  -> **NOTE:** If `ipv6_enable` is set to **true**, the subnet must support **IPV6**.

* `enable_all` - (Required, Bool, ForceNew) Specifies whether a bucket instance is a family bucket instance.
  - **true**. Create a Family Bucket instance.
  - **false**. Create a MQS instance.

  Changing this will create a new instance.

* `cpu_architecture` - (Required, String, ForceNew) Specifies the architecture type of instance, the valid value are
  as follows.
  - **x86**. Create a x86 architecture instance.
  - **arm**. Create a arm architecture instance.

  Changing this will create a new instance.

* `mqs` - (Required, Map, ForceNew) Specifies MQS service parameters.

  The [mqs](#roma_instance_mqs_arg) structure is documented below.

  Changing this will create a new instance.

* `maintain_begin` - (Optional, String, ForceNew) Maintain start time. Example **22:00:00**.

  Changing this will create a new instance.

* `maintain_end` - (Optional, String, ForceNew) Maintain end time. Example **02:00:00**.

  Changing this will create a new instance.

-> **NOTE:** Parameters `maintain_begin` and `maintain_end` must be set in pairs.

~> **WARNING:** The `maintain_begin` and `maintain_end` will be **Deprecated** in later version.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID.

  Changing this will create a new instance.

* `entrance_bandwidth_size` - (Optional, Int, ForceNew) Specifies the ingress bandwidth. This parameter is **Required**
  when `enable_publicip` is **true**.

  Changing this will create a new instance.

<a name="roma_instance_mqs_arg"></a>
The `mqs` block supports:

* `enable_publicip` - (Optional, Bool, ForceNew) Specifies whether public network access is supported.
  Default to **false**.

  Changing this will create a new instance.

* `ssl_enable` - (Optional, Bool, ForceNew) Specifies whether public network access is supported.
  Default to **false**.

  Changing this will create a new instance.

* `engine_version` - (Optional, String, ForceNew) Specifies the **Kafka** engine version. This parameter is
  **conflicted** with `rocketmq_enable`.  
  The valid values are as follows:
  + **2.7**
  + **2.3.0**
  + **1.1.0**

  Changing this will create a new instance.

* `rocketmq_enable` - (Optional, Bool, ForceNew) Specifies create a **RocketMQ** instance. This parameter is
  **conflicted** with `engine_version`. Default to **false**.

  Changing this will create a new instance.

* `enable_acl` - (Optional, Bool, ForceNew) Specifies whether to enable acl. This parameter is valid only when the
  `rocketmq_enable` is **true**.

  Changing this will create a new instance.

* `retention_policy` - (Optional, String, ForceNew) Specifies the message integration capacity threshold policy.
  This parameter is valid only when the engine type is **Kafka**.
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

* `site_id` - The ID of a site instance. This parameter is valid only in the site scenario.

* `flavor_id` - The instance flavor ID.

* `flavor_type` - The instance flavor type. The valid values are as follows:
  + **basic**
  + **professional**
  + **enterprise** 
  + **platinum**
  + **platinumXB**

* `cpu_arch` - The CPU architecture type. The valid values are as follows:
  - **x86_64**: x86 architecture. 
  - **aarch64**: ARM architecture.

* `publicip_id` - The EIP ID bound to the instance.

* `publicip_address` -The EIP address bound to the instance.

* `publicip_enable` - Whether to enable public network access. When public network access is enabled,
  `publicip_id` is **Required**.

* `connect_address` - The ROMA Connect connection address.

* `charge_type` - The charging type of the resource specification.
  + **postPaid**

* `bandwidths` - The egress bandwidth of the public network.

* `ipv6_connect_address` - The ROMA Connect connection address in IPv6 scenarios.

* `external_elb_enable` - Whether to enable the external ELB.

* `external_elb_id` - The IP address ID of the external ELB.

* `external_elb_address` - The IP address of the external ELB.

* `external_eip_bound` - Whether the external ELB is bound to a public IP address.

* `external_eip_id` - The ID of the public network address bound to the external load balancer.

* `external_eip_address` - The public IP address bound to the external ELB.

* `create_time` - The creation time.

* `update_time` - The last update time.

* `resources` - The description of resource information.

  The [resources](#roma_instance_resources_attr) structure is documented below.

* <a name="roma_instance_resources_attr"></a>
The `resources` block supports:

* `mqs` - The MQS service parameters. The [mqs](#roma_instance_mqs_attr) structure is documented below.

<a name="roma_instance_mqs_attr"></a>
The `mqs` block supports:

* `id` - The instance ID of ROMA Connect.

* `enable` - Whether to enable the MQS component.

* `retention_policy` - The message integration capacity threshold policy.

* `ssl_enable` - Whether to enable SSL.

* `trace_enable` - Whether to enable the message track.

* `vpc_client_plain` - Whether to enable plaintext access in a VPC.

* `partition_num` - The number of partitions.

* `specification` - The specification of instance.

* `private_connect_address` - The internal connection address.

* `public_connect_address` - The external connection address.

* `private_restful_address` - The internal restful connection address.

* `public_restful_address` - The external restful connection address.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.

* `delete` - Default is 30 minutes.
