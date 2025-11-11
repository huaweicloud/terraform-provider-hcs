---
subcategory: CodeArts
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_codearts_project"
description: |-
  Manages a Project resource within HuaweiCloudStack.
---

# hcs_codearts_project

-> **NOTE:** This resource can only be used in HCS **8.5.0** and **later** version.

Manages a Project resource within HuaweiCloudStack.

## Example Usage

```hcl
resource "hcs_codearts_project" "test" {
  name = "demo_project"
  type = "scrum"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) The project name.  
  The name can contain `1` to `128` characters.

* `type` - (Required, String, ForceNew) The type of project. The valid values are as follows:
  - **scrum**
  - **xboard**
  - **basic**
  - **phoenix**
  - **ipd**

  Changing this parameter will create a new resource.

* `description` - (Optional, String) The description about the project.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project ID of the project.  
  Value 0 indicates the default enterprise project.

  Changing this parameter will create a new resource.

* `source` - (Optional, String, ForceNew) The source of project.

  Changing this parameter will create a new resource.

* `template_id` - (Optional, Int, ForceNew) The template id which used to create project. This parameter is **Required**
  when `type` is **ipd**.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `archive` - Whether the project is archived.

* `project_code` - The project code.

* `project_num_id` - The number id of project.

## Import

The project can be imported using the `id`, e.g.

```bash
$ terraform import hcs_codearts_project.test 0ce123456a00f2591fabc00385ff1234
```
