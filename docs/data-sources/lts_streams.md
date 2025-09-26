---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_lts_streams"
description: |-
  Use this data source to get the list of LTS log streams.
---

# hcs_lts_streams

Use this data source to get the list of LTS log streams.

## Example Usage

```hcl
data "hcs_lts_streams" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the log streams.  
  If omitted, the provider-level region will be used.

* `log_group_name` - (Optional, String) Specifies the name of the log group to be queried.

* `log_stream_name` - (Optional, String) Specifies the name of the log stream to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `log_streams` - The list of log streams.

  The [log_streams](#lts_log_stream_attr) structure is documented below.

<a name="lts_log_stream_attr"></a>
The `streams` block supports:

* `id` - The ID of the log stream.

* `name` - The name of the log stream.

* `created_at` - The creation time of the log stream, in RFC3339 format.

* `filter_count` - The number of filters.

* `tags` - The tags of the log stream.

* `is_favorite` - Whether the stream is bookmarked.
