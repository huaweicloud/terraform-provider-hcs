---
subcategory: "Virtual Data Center (VDC)"
---

# hcs_vdc_role

Use this data source to get details of the specified VDC role.

## Example Usage

```hcl

data "hcs_vdc_role" "test" {
  display_name = "Tenant Guest"
}

```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional, String)  Specifies the display name of the role displayed on the console. It is
  recommended to use this parameter instead of name and required if `name` is not specified.

* `name` - (Optional, String) Specifies the name of the role for internal use. It's required if `display_name` is not
  specified. The name must meet the following requirements:

    * If you query a system policy, the name can be `system_all_xxx`.
    * If you query a system role, the name can be `te_agency` or `te_admin`.
    * If you query a custom role, the name can be `custom_xxx`.

* `role_type` - (Optional, String) Role type. Valid options are as follows:

    * system: System-defined role and System-defined policy
    * custom: Custom role

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the role.
* `description` - The description of the role.
* `type` - The display mode of the role.
* `policy` - The content of the role.
* `catalog` - The service catalog of the role.
