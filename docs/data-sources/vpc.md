---
subcategory: "Virtual Private Cloud (VPC)"
---

# hcs_vpc

Provides details about a specific VPC.

## Example Usage

```hcl
variable "vpc_name" {}

data "hcs_vpc" "vpc" {
  name = var.vpc_name
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available VPCs in the current region. The given
filters must match exactly one VPC whose data will be exported as attributes.

* `region` - (Optional, String) Specifies the region in which to obtain the VPC. If omitted, the provider-level region
  will be used.

* `name` - (Optional, String) Specifies an unique name for the VPC. The value is a string of no more than 64 characters
  and can contain digits, letters, underscores (_), and hyphens (-).

* `id` - (Optional, String) Specifies the id of the VPC to retrieve.

* `status` - (Optional, String) Specifies the current status of the desired VPC. The value can be CREATING, OK or ERROR.

* `cidr` - (Optional, String) Specifies the cidr block of the desired VPC.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID which the desired VPC belongs to.

* `secondary_cidrs` - The secondary CIDR blocks of the VPC.

## Attributes Reference

