---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_lts_structuring_configuration"
description: ""
---

# hcs_lts_structuring_configuration

Manages an LTS structuring configuration resource within HuaweiCloudStack.

## Example Usage

### Creating with system template

```hcl
variable "log_group_id" {}
variable "log_stream_id" {}

resource "hcs_lts_structuring_configuration" "test" {
  log_group_id  = var.log_group_id
  log_stream_id = var.log_stream_id
  template_name = "CTS"
  template_type = "built_in"
}
```

### Creating with custom template

```hcl
variable "log_group_id" {}
variable "log_stream_id" {}
variable "template_name" {}
variable "template_id" {}

resource "hcs_lts_structuring_configuration" "test" {
  log_group_id  = var.log_group_id
  log_stream_id = var.log_stream_id
  template_name = var.template_name
  template_id   = var.template_id
  template_type = "custom"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `log_group_id` - (Required, String, ForceNew) Specifies the log group ID.
  Changing this parameter will create a new resource.

* `log_stream_id` - (Required, String, ForceNew) Specifies the log stream ID.
  Changing this parameter will create a new resource.

* `template_type` - (Required, String) Specifies the type of the template. The valid values are as follows:
  + **built_in**: System templates.
  + **custom**:   Custom templates.

* `template_name` - (Required, String) Specifies the template name. When `template_type` is set to **built_in**,
  valid values are:
  + **ELB**
  + **VPC**
  + **DCS_AUDIT**: DCS audit log.
  + **TOMCAT**
  + **NGINX**

* `template_id` - (Optional, String) Specifies the template ID. The field is valid and required only when
  `template_type` is set to **custom**.

* `demo_fields` - (Optional, List) Specifies the example fields. Use to set quick analysis configurations for fields.
  Only need to enter the fields whose status is different from that of `is_analysis` in the template.
  The [demo_fields](#StructConfig_fields) structure is documented below.

* `tag_fields` - (Optional, List) Specifies the tag fields. Use to set quick analysis configurations for fields.
  Only need to enter the fields whose status is different from that of `is_analysis` in the template.
  The [tag_fields](#StructConfig_fields) structure is documented below.

* `quick_analysis` - (Optional, Bool) Specifies whether to enable `demo_fields` and `tag_fields` quick analysis.
  + If this parameter is set to **true**, quick analysis is enabled for all `demo_fields` and `tag_fields`.
  + If this parameter is set to **false**, `is_analysis` in `demo_fields` and `tag_fields` in the template determines
    whether to enable quick analysis.

  Defaults to **false**.

<a name="StructConfig_fields"></a>
The `demo_fields` and `tag_fields` block supports:

* `is_analysis` - (Optional, Bool) Specifies whether quick analysis is enabled. Defaults to **false**.

* `field_name` - (Required, String) Specifies the field name. The valid length is limited from `1` to `64`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `demo_log` - The sample log event.

## Import

The LTS structuring configuration can be imported using `log_group_id` and `log_stream_id`, separated by a slash, e.g.

```bash
$ terraform import hcs_lts_structuring_configuration.test <log_group_id>/<log_stream_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `template_type`, `template_id`,
`demo_fields`, `tag_fields`, `quick_analysis`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "hcs_lts_structuring_configuration" "test" {
  ...
  
  lifecycle {
    ignore_changes = [
      template_type,
      template_id,
      demo_fields,
      tag_fields,
      quick_analysis,
    ]
  }
}
```
