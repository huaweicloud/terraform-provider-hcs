---
subcategory: "Enterprise Projects (EPS)"
---

# hcs_enterprise_project

Manages an EPS resource within HuaweiCloudStack.

## Example Usage

```hcl
resource "hcs_enterprise_project" "eps1" {
  name = "test_eps1"
  project_id = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  description = "test_eps1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The value can contain 1 to 64 characters consisting of letters, digits, underscores (_), and hyphens (-). The name cannot be changed to default (case insensitive), and must be unique in the domain.

* `project_id` - (Required, String) Resource set. The value can contain 1 to 36 characters, including only lowercase letters, digits, and hyphens (-).

* `description` - (Optional, String) Description. The value can contain 0 to 512 characters excluding angle brackets (< and >).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Enterprise project ID.

* `name` - (Required, String) The value can contain 1 to 64 characters consisting of letters, digits, underscores (_), and hyphens (-). The name cannot be changed to default (case insensitive), and must be unique in the domain.

* `project_id` - (Required, String) Resource set. The value can contain 1 to 36 characters, including only lowercase letters, digits, and hyphens (-).

* `description` - (Optional, String) Description. The value can contain 0 to 512 characters excluding angle brackets (< and >).

## Import

Instances can be imported by their `id`. For example,

```
terraform import hcs_enterprise_project.eps8 267f2c1f-8ff5-4443-a067-439d72107a29
```

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minute.
* `update` - Default is 5 minute.
* `delete` - Default is 5 minute.