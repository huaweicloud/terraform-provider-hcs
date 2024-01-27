---
subcategory: "Elastic Volume Service (EVS)"
---

# hcs_evs_volume_types

Use this data source to query the detailed information list of the volume types within HuaweiCloudStack.

## Example Usage

```hcl
variable "availability_zone" {}

data "hcs_evs_volume_types" "test" {
  availability_zone = var.availability_zone
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional, String) Specifies the availability zone for the volume type.

* `extra_specs` - (Optional, Map) Extra attribute key/value pairs to associate with the volume type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A data source ID in hashcode format.

* `volume_types` - The detailed information of the volume type. Structure is documented below.

The `volume_types` block supports:

* `id` - The resource ID of volume type, in UUID format.

* `name` - The volume type name.

* `is_public` - Whether the volume type is publicly visible.

* `description` - The volume type description.

* `extra_specs` - Extra attribute key/value pairs to associate with the volume type.

* `qos_specs_id` - The ID of qos.
