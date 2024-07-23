---
subcategory: "Virtual Private Cloud (VPC)"
---

# hcs_vpc

Manages a VPC resource within HuaweiCloudStack.

## Example Usage

```hcl
variable "vpc_name" {
  default = "hcs_vpc"
}

variable "vpc_cidr" {
  default = "192.168.0.0/16"
}

variable "vpc_secondary_cidrs" {
  default = ["192.182.0.0/24","192.183.0.0/24","192.184.0.0/24","192.185.0.0/24"]
}

resource "hcs_vpc" "vpc" {
  name = var.vpc_name
  cidr = var.vpc_cidr
  secondary_cidrs = var.vpc_secondary_cidrs
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the VPC. If omitted, the
  provider-level region will be used. Changing this creates a new VPC resource.

* `name` - (Required, String) Specifies the name of the VPC. The name must be unique for a tenant. The value is a string
  of no more than 64 characters and can contain digits, letters, underscores (_), and hyphens (-).

* `cidr` - (Required, String) Specifies the range of available subnets in the VPC. The value ranges from 10.0.0.0/8 to
  10.255.255.0/24, 172.16.0.0/12 to 172.31.255.0/24, or 192.168.0.0/16 to 192.168.255.0/24.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID which the desired VPC belongs to.

* `secondary_cidrs` - (Optional, Set) Specifies the secondary CIDR blocks of the VPC.
  Each VPC can have 4 secondary CIDR blocks.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The VPC ID in UUID format.

* `status` - The current status of the VPC. Possible values are as follows: CREATING, OK or ERROR.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `delete` - Default is 3 minute.

## Import

VPCs can be imported using the `id`, e.g.

```
$ terraform import hcs_vpc.vpc_v1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
