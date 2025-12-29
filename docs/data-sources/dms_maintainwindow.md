---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_dms_maintainwindow"
description: |-
  Use this data source to get the ID of an available HuaweiCloudStack dms maintenance windows.
---

# hcs_dms_maintainwindow

Use this data source to get the ID of an available HuaweiCloudStack dms maintenance windows.

## Example Usage

```hcl
data "hcs_dms_maintainwindow" "test" {
  seq = 1
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the dms maintenance windows. If omitted, 
  the provider-level region will be used.

* `seq` - (Optional, Int) Specifies the sequential number of a maintenance time window.

* `begin` - (Optional, String) Specifies the time at which a maintenance time window starts.

* `end` - (Optional, String) Specifies the time at which a maintenance time window ends.

* `default` - (Optional, Bool) Whether a maintenance time window is set to the default time segment.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates a data source ID in UUID format.
