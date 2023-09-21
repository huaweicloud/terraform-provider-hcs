---
subcategory: "Enterprise Projects (EPS)"
---

# hcs_enterprise_project

Use this data source to get the list of the enterprise projects.

## Example Usage

```hcl
data "hcs_enterprise_project" "epslist" {
    name = "eps1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Enterprise project name. Fuzzy matching is supported.

* `id` - (Optional, String) Enterprise project ID. Value 0 indicates that the default enterprise project is queried.

* `ids` - (Optional, String) List of enterprise project IDs, separated by commas (,). Up to 100 enterprise project IDs are supported.

* `domain_id` - (Optional, String) Tenant ID. This parameter is mandatory for the op_service permission.

* `vdc_id` - (Optional, String) VDC ID, used to query the enterprise project list in a specified VDC. This parameter can only be used by users who have the permission to manage enterprise projects.

* `inherit` - (Optional, Bool) Whether to query enterprise projects in lower-level VDCs. This parameter is valid only when the VDC_ID parameter has been set.

* `project_id` - (Optional, String) Resource set ID, used to query the enterprise project list in a specified resource space.

* `type` - (Optional, String) Enterprise project type. The value can be prod or poc. Currently, the poc type is not supported. The default value is prod.

* `status` - (Optional, Int) Enterprise project status. The value can be 1 or 2. 1 indicates that the enterprise project is enabled, and 2 indicates that the enterprise project is disabled. Currently, a enterprise project cannot be disabled. The default value is 1.

* `query_type` - (Optional, String) Query type. The value can be auth, auths, or list. The default value is auth. Value auth indicates that the enterprise projects on which the user has permissions is queried in a resource space where the token must be obtained. Value auths indicates that all enterprise projects on which the user has permissions are queried. Value list indicates that all enterprise projects are queried.

* `auth_action` - (Optional, String) Enterprise projects on which you have a specified permission. This parameter is available only when query_type is set to auth.

* `contain_default` - (Optional, Bool) Whether the default enterprise project is included. The statistics page on the administrator portal does not count the default enterprise project.

* `offset` - (Optional, String) Index position, which starts from the next data record specified by offset. The value must be a number and cannot be a negative number. The default value is 0. The maximum value of offset is 1000000000.

* `limit` - (Optional, String) The number of records to be queried, which ranges from 1 to 1000. The default value is 1000.

* `sort_key` - (Optional, String) Sorting field. The value can be created_at or updated_at. The default value is created_at.

* `sort_dir` - (Optional, String) Sorting direction. The value can be asc or desc. The default value is desc.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - List of EPS instance details. The object structure of each EPS instance is documented below.

The `instances` block supports:

* `id` - Enterprise project ID.

* `name` - Enterprise project name.

* `description` - Description.

* `type` - Type.

* `delete_flag` - Whether the enterprise project can be deleted.

* `status` - Status. Value 1 indicates the enterprise project is enabled. Value 2 indicates the enterprise project is disabled.

* `created_at` - Creation time.

* `updated_at` - Modification time.

* `domain_id` - ID of a Tenant to which the enterprise project belongs.

* `vdc_id` - ID of a VDC to which the enterprise project belongs.

* `project_id` - ID of a resource space to which the enterprise project belongs.

* `domain_name` - Name of a tenant to which the enterprise project belongs.

* `vdc_name` - Name of a VDC to which the enterprise project belongs.

* `project_name` - Name of a resource space to which the enterprise project belongs.