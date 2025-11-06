---
subcategory: "GaussDB"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_gaussdb_opengauss_parameter_template"
description: |-
  Use this data source to get the details of a GaussDB configuration.
---

# hcs_gaussdb_opengauss_parameter_template

Use this data source to get the details of a GaussDB configuration.

## Example Usage

```hcl
variable "template_id" {}

data "hcs_gaussdb_opengauss_parameter_template" "test" {
  template_id = var.template_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the instance tags.

* `template_id` - (Required, String) Specifies the parameter template ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID.

* `name` - Indicates the name of the configuration.

* `description` - Indicates the description of the configuration.

* `engine_version` - Indicates the engine version.

* `instance_mode` - Indicates the instance mode.

* `created_at` - Indicates the creation time.

* `updated_at` - Indicates the update time.

* `parameters` - Indicates the list of configuration parameters.
  The [parameters](#configuration_parameters) structure is documented below.

<a name="configuration_parameters"></a>
The `parameters` block supports:

* `name` - Indicates the name of the parameter.

* `value` - Indicates the value of the parameter.

* `need_restart` - Indicates whether the parameter needs to be restarted.

* `readonly` - Indicates whether the parameter is read-only.

* `value_range` - Indicates the value range of the parameter.

* `data_type` - Indicates the data type of the parameter.

* `description` - Indicates the description of the parameter.

* `is_risk_parameter` - Whether this parameter is a risk parameter.

* `risk_description` - Indicates the description of risk parameter.
