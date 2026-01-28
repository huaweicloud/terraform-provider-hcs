---
subcategory: "MapReduce Service (MRS)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_mrs_versions"
description: |-
  Use this data source to get available cluster versions of MapReduce.
---

# hcs_mrs_versions

Use this data source to get available cluster versions of MapReduce.

## Example Usage

```hcl
data "hcs_mrs_versions" "test" {
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `versions` - List of available cluster versions.
