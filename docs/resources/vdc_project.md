---
subcategory: "Virtual Data Center (VDC)"
---

# hcs_vdc_project

Manages a VDC resource space within Huawei Cloud Stack.

## Example Usage

```hcl

variable "vdc_id" {}

resource "hcs_vdc_project" "test" {
  vdc_id = var.vdc_id
  name   = "cn-global-1_project1"
}

```

## Argument Reference

The following arguments are supported:

* `vdc_id` - (Required, String) ID of the VDC that the created resource space belongs to. The ID can contain 1 to 36
  characters. The VDC ID cannot be changed.

* `name` - (Required, String) Indicates the resource space name. The resource space name must start with `{region_id}_`
  and can contain only letters (case-insensitive), digits, hyphens (-), underscores (_), and parentheses. The name can
  contain 1 to 64 characters.

* `display_name` - (Optional, String) Indicates the display name of a resource space. If this parameter is not
  transferred, the system automatically generates a display name based on the resource space name. The display name can
  contain 0 to 64 characters and cannot contain greater-than signs (>) or less-than signs (<).

* `description` - (Optional, String) Indicates the description of a resource space. The description can contain 0 to 255
  characters and cannot contain greater-than signs (>) or less-than signs (<).。

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource space ID.

* `regions` - The region list.

## import

VDC resource spaces can be imported using the `id`, e.g.

```bash
$ terraform import hcs_vdc_project.project1 0350a018560a499692d972749fa6c94c
```
