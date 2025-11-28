---
subcategory: "ServiceStage"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_servicestage_environment"
description: |-
  Manages a ServiceStage environment resource within HuaweiCloudStack
---

# hcs_servicestage_environment

Manages a ServiceStage environment resource within HuaweiCloudStack.

## Example Usage

### Create an environment based on some cce clusters

```hcl
variable "env_name" {}
variable "vpc_id" {}
variable "cce_cluster_id" {}
variable "cci_namespace_name" {}

resource "hcs_servicestage_environment" "test" {
  name   = var.env_name
  vpc_id = var.vpc_id

  basic_resources {
    type = "cce"
    id   = var.cce_cluster_id
  }
  basic_resources {
    type = "cci"
    id   = var.cci_namespace_name
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the ServiceStage environment.
  If omitted, the provider-level region will be used.

  Changing this will create a new resource.

* `name` - (Required, String) Specifies the environment name.
  The name can contain of `2` to `64` characters, only letters, digits, hyphens (-) and underscores (_) are allowed.
  The name must start with a letter and end with a letter or a digit.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID to which the environment belongs.

  Changing this will create a new resource.

* `deploy_mode` - (Optional, String, ForceNew) Specifies the environment type. The valid values ars as follows:
  + **virtualmachine**: Virtual machine type
  + **container**: Kubernetes type
  + **mixed**: Virtual machine and kubernetes type

  Changing this will create a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the application
  belongs.

  Changing this will create a new resource.

* `basic_resources` - (Required, List) Specifies the basic resources.
  The [basic_resources](#servicestage_env_resources) object structure is documented below.

* `optional_resources` - (Optional, List) Specifies the optional resources.
  The [optional_resources](#servicestage_env_resources) object structure is documented below.

* `description` - (Optional, String) Specifies the environment description.
  The description can contain a maximum of `128` characters.

* `alias` - (Optional, String, ForceNew) Specifies the environment alias.
  The alias can contain a maximum of `64` characters.

  Changing this will create a new resource.

<a name="servicestage_env_resources"></a>
The `basic_resources` and `optional_resources` block supports:

* `type` - (Required, String) Specifies the resource type.
  + The type of `basic` resource supports **cce**, **ecs**.
  + The type of `optional` resource supports **rds**, **dcs**, **elb**, **cse**, **nacos** and **gateway**.

* `id` - (Required, String) Specifies the resource ID. For most resources, this parameter needs to fill in their **id**,
  but for CCI namespace, this parameter needs to fill in **name**.

-> All resources must under the same VPC as the environment.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The environment ID in UUID format.

## Import

Environments can be imported using their `id`, e.g.:

```bash
$ terraform import hcs_servicestage_environment.test 17383329-b686-47e4-8f70-0d8dcddb65e9
```
