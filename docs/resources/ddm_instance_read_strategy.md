---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_ddm_instance_read_strategy"
description: |-
  Use this resource to set a DDM instance read strategy within HuaweiCloudStack.
---

# hcs_ddm_instance_read_strategy

Use this resource to set a DDM instance read strategy within HuaweiCloudStack.

-> **NOTE:** This resource is a one-time action resource using to read ddm instance strategy. Deleting this resource
will not clear the corresponding request record, but will only remove the resource information from the tfstate file

## Example Usage

```hcl
variable "instance_id" {}
variable "rds_instance_id" {}
variable "rds_read_replica_instance_id" {}

resource "hcs_ddm_instance_read_strategy" "test" {
  instance_id = var.instance_id

  read_weights {
    db_id  = var.rds_instance_id
    weight = 70
  }

  read_weights {
    db_id  = var.rds_read_replica_instance_id
    weight = 30
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DDM instance.

* `read_weights` - (Required, List) Specifies the list of read weights of the primary DB instance
  and its read replicas. The valid value is range from `0` to `100`.

  The [read_weights](#ddm_instance_read_strategy_read_weights) object structure is documented below.

<a name="ddm_instance_read_strategy_read_weights"></a>
The `read_weights` block supports:

* `db_id` - (Required, String) Specifies the ID of the DB instance associated with the DDM schema.

* `weight` - (Required, Int) Specifies read weight of the DB instance associated with the DDM schema.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID. The value is the DDM instance ID.
