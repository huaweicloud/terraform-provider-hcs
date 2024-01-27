---
subcategory: "Web Application Firewall (WAF)"
---

# hcs_waf_reference_table

Manages a WAF reference table resource within HuaweiCloudStack.

-> **NOTE:** All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be
used. The reference table resource can be used in Cloud Mode (professional version), Dedicated Mode and ELB Mode.

## Example Usage

```hcl
variable "enterprise_project_id" {}

resource "hcs_waf_reference_table" "ref_table" {
  name                  = "tf_ref_table_demo"
  type                  = "url"
  enterprise_project_id = var.enterprise_project_id

  conditions = [
    "/admin",
    "/manage"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the WAF reference table resource. If omitted,
  the provider-level region will be used. Changing this setting will push a new reference table.

* `name` - (Required, String) The name of the reference table. Only letters, digits, and underscores(_) are allowed. The
  maximum length is 64 characters.

* `type` - (Required, String, ForceNew) The type of the reference table, The options are `url`, `user-agent`,
  `params`, `cookie`, `referer` and `header`. Changing this setting will push a new reference table.

* `conditions` - (Required, List) The conditions of the reference table. The maximum length is 30. The maximum length of
  condition is 2048 characters.

* `description` - (Optional, String) The description of the reference table. The maximum length is 128 characters.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID of WAF reference table.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the reference table.

* `creation_time` - The server time when reference table was created.

## Import

There are two ways to import WAF reference table state.

* Using the `id`, e.g.

```bash
$ terraform import hcs_waf_reference_table.test <id>
```

* Using `id` and `enterprise_project_id`, separated by a slash, e.g.

```bash
$ terraform import hcs_waf_reference_table.test <id>/<enterprise_project_id>
```
