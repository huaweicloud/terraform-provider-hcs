---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_cfw_protection_rule_hit_count"
description: ""
---

# hcs_cfw_protection_rule_hit_count

Use this data source to get cfw protection rule hit count.

-> **NOTE:** To use this data source, your HCS version must be **8.5.0** or later.

## Example Usage

```hcl
variable "rule_id" {}

data "hcs_cfw_protection_rule_hit_count" "test" {
  rule_ids = [var.rule_id]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `rule_ids` - (Required, List) Specifies the protection rule id list.

* `fw_instance_id` - (Optional, String) Specifies the firewall instance ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `records` - The firewall instance records.
  The [records object](#records) structure is documented below.

<a name="records"></a>
The `records` block supports:

* `rule_id` - The rule id.

* `rule_hit_count` - Number of times that an ACL rule is hit. When an ACL rule is triggered,
  the number of times that the rule ID is hit is added.
