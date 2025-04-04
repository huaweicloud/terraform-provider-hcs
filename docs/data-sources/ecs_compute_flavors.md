---
subcategory: "Elastic Cloud Server (ECS)"
---

# hcs_ecs_compute_flavors

Use this data source to get the available Compute Flavors.

## Example Usage

```hcl
data "hcs_ecs_compute_flavors" "flavors" {
  availability_zone = "az1.dc1"
  cpu_core_count    = 2
  memory_size       = 4
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the flavors.
  If omitted, the provider-level region will be used.

* `availability_zone` - (Optional, String) Specifies the AZ name.

* `cpu_core_count` - (Optional, Int) Specifies the number of vCPUs in the ECS flavor.

* `memory_size` - (Optional, Int) Specifies the memory size(GB) in the ECS flavor.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates a data source ID.

* `ids` - A list of flavor IDs.

* `flavors` - List of ECS flavors details. The [flavor object](#compute_instance_flavor_object) structure of each flavor is documented below.

<a name="compute_instance_flavor_object"></a>
The `flavor` block supports:

* `id` - The flavor ID.
* `name` - The flavor Name.
* `ram` - Memory size specified by the flavor.
* `vcpus` - Number of virtual CPU cores specified by the flavor.
* `ext_boot_type` - The system disk type, which can be a *Volume* or *LocalDisk*.
