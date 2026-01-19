---
subcategory: "Virtual Data Center (VDC)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_vdc_agency"
description: ""
---

# hcs_vdc_agency

Manages a VDC Agency within Huawei Cloud Stack.

-> **NOTE:** Supported from ManageOne version 8.6.1 onwards.

## Example Usage

```hcl
resource "hcs_vdc_agency" "agency" {
  name                  = "test_agency"
  delegated_domain_name = "test-vdc"
  project_role {
    project = "project-1"
    roles   = ["VDC Admin"]
  }
  domain_roles        = ["VDC ReadOnly"]
  all_resources_roles = ["Tag Admin"]
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew): Agency name, which can contain 1 to 64 characters. Only letters, digits, hyphens (-), underscores (_), and periods (.) are allowed.

* `description` - (Optional, String, ForceNew) Description. The value cannot contain the following characters: >< The value can
  contain 0 to 255 characters.

* `delegated_domain_name` - (Required, String, ForceNew): Name of an existing maintained tenant.

* `project_role` - (Optional, List): Permissions that can be granted for accessing resource spaces. You can use multiple project_role blocks authorize different resource spaces. Do not use the same resource space name in different project_role blocks. The project_role object structure is documented below.

* `domain_roles` - (Optional, List): Optional permissions that can be granted for accessing global services. One type of such permission can be granted to multiple tenants.

* `all_resources_roles` - (Optional, List): Permissions that can be granted for accessing global resources. One type of such permission can be granted to multiple tenants. These permission settings can be inherited by resource spaces.

The `project_role` block supports:

* `project` - (Required, String): Name of the resource space for which the permissions are granted. The resource space must belong to the current tenant.

* `roles` - (Required, List): Mandatory permissions that can be granted for accessing resource spaces. One type of such permission can be grated at a time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Agency ID.

## Import

Agencies can be imported using the id, e.g.

```bash
$ terraform import hcs_vdc_agency.agency 93ea9bae6e314377bfc22078294e18ac
```
