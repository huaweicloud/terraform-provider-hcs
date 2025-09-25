---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_lts_groups"
description: |-
  Use this data source to get the list of LTS log groups.
---

# hcs_lts_groups

Use this data source to get the list of LTS log groups.

## Example Usage

```hcl
data "hcs_lts_groups" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the log groups.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `log_groups` - The list of log groups.

  The [log_groups](#lts_log_groups_attr) structure is documented below.

<a name="lts_log_groups_attr"></a>
The `log_groups` block supports:

* `id` - The ID of the log group.

* `name` - The name of the log group.

* `created_at` - The creation time of the log group, in RFC3339 format.

* `ttl_in_days` - The storage duration of the log group in days.

* `stream_size` - The number of log streams under the log group.

* `tags` - The tags of the log group.
