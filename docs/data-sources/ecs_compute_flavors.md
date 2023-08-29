---
subcategory: "Elastic Cloud Server (ECS)"
---

# hcs_ecs_compute_flavors

Use this data source to get the available Compute Flavors.

## Example Usage

```hcl
data "hcs_ecs_compute_flavors" "flavors" {
  availability_zone = "az1.dc1"
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the flavors.
  If omitted, the provider-level region will be used.

* `availability_zone` - (Optional, String) Specifies the AZ name.

* `performance_type` - (Optional, String) Specifies the ECS flavor type. Possible values are as follows:
  + **normal**: General computing
  + **computingv3**: General computing-plus
  + **highmem**: Memory-optimized
  + **saphana**: Large-memory HANA ECS
  + **diskintensive**: Disk-intensive

* `generation` - (Optional, String) Specifies the generation of an ECS type. For example, **s3** indicates
  the general-purpose third-generation ECSs. For details, see
  [ECS Specifications](https://support.huaweicloud.com/intl/en-us/productdesc-ecs/ecs_01_0014.html).

* `cpu_core_count` - (Optional, Int) Specifies the number of vCPUs in the ECS flavor.

* `memory_size` - (Optional, Int) Specifies the memory size(GB) in the ECS flavor.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates a data source ID.

* `ids` - A list of flavor IDs.
