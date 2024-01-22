---
subcategory: "Elastic Volume Service (EVS)"
---

# hcs_evs_snapshots

Use this data source to query the detailed information list of the EVS snapshots within HuaweiCloudStack.

## Example Usage

```hcl
variable "volume_id" {}
data "hcs_evs_snapshots" "snapshots" {
  volume_id = var.volume_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the snapshot list.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specify the enterprise project ID to filter.
  If this field value is not specified, snapshots of all enterprise projects within
  authority scope will be queried.

* `name` - (Optional, String)  The name of snapshot. Maximum supported is 64 characters.

* `status` - (Optional, String) The status of snapshot.

* `volume_id` - (Optional, String) The ID of the disk to which the snapshot belongs.

* `availability_zone` - (Optional, String) The availability zone of the disk to which the snapshot belongs.

* `snapshot_id` - (Optional, String) Specify the snapshot ID to filter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, in UUID format.

* `snapshots` - A list of EVS snapshots.

The `snapshots` block supports:

* `id` - The resource ID of EVS snapshot.

* `name` - The name of snapshot. Maximum supported is 64 characters.

* `size` - The snapshot size. Unit is GiB.

* `status` - The status of snapshot.

* `description` - The snapshot description information.

* `created_at` - The snapshot creation time.

* `updated_at` - The snapshot update time.

* `volume_id` - The ID of the disk to which the snapshot belongs.