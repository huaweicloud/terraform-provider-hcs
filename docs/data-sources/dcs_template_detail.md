---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_dcs_template_detail"
description: |-
  Use this data source to get the detail of DCS template.
---

# hcs_dcs_template_detail

Use this data source to get the detail of DCS template.

## Example Usage

```hcl
variable "template_id" {}

data "hcs_dcs_template_detail" "test" {
  template_id = var.template_id
  type        = "sys"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `type` - (Required, String) Specifies the type of the template. The valid values are as follows:
  + **sys**: system template.
  + **user**: custom template.

* `template_id` - (Required, String) Specifies the ID of the template.

* `params` - (Optional, List) Specifies the list of the template params.  
  The [params](#rds_params_arg) object structure is documented below.

<a name="rds_params_arg"></a>
The `params` block supports:

* `param_name` - (Optional, String) Specifies the name of the param.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `name` - Indicates the name of the template.

* `type` - Indicates the type of the template. The valid values are as follows:
  + **sys**
  + **user**

* `engine` - Indicates the cache engine. Currently, only **Redis** is supported.

* `engine_version` - Indicates the cache engine version. The valid values are as follows:
  + **3.0**
  + **4.0**
  + **5.0**
  + **6.0**.

* `cache_mode` - Indicates the DCS instance type. The valid values are as follows: 
  + **single**
  + **ha**
  + **cluster**
  + **proxy**
  + **ha_rw_split**

* `product_type` - Indicates the product edition. The valid values are as follows:
  + **generic**
  + **enterprise**

* `storage_type` - Indicates the storage type. The valid values are as follows:
  + **DRAM**
  + **SSD**

* `description` - Indicates the description of the template.

* `params` - Indicates the list of the template params.  
  The [params](#rds_params_attr) object structure is documented below.

<a name="rds_params_attr"></a>
The `params` block supports:

* `param_id` - Indicates the ID of the param.

* `default_value` - Indicates the default of the param.

* `value_range` - Indicates the value range of the param.

* `value_type` - Indicates the value type of the param.

* `description` - Indicates the description of the param.

* `need_restart` - Indicates whether the DCS instance need restart.
