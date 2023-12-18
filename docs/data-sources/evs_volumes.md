---
subcategory: "Elastic Volume Service (EVS)"
---

# hcs_evs_volumes

Use this data source to query the detailed information list of the EVS disks within HuaweiCloudStack.

## Example Usage

```hcl
data "hcs_evs_volumes" "test" {
  name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the disk list.
  If omitted, the provider-level region will be used.

* `availability_zone` - (Optional, String) Specifies the availability zone for the disks.

* `multiattach` - (Optional, Bool) Specifies whether the disk is shareable.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `metadata` - (Optional, Map) Metadata key/value pairs to associate with the disk.

* `name` - (Optional, String) The disk name.

* `status` - (Optional, String) Specifies the disk status.

* `tags` - (Optional, Map) The key/value pairs to associate with the disk.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A data source ID in hashcode format.

* `volumes` - The detailed information of the disks. Structure is documented below.

The `volumes` block supports:

* `id` - The resource ID of EVS disk, in UUID format.

* `attachments` - The disk attachment information. Structure is documented below.

* `availability_zone` - The availability zone of the disk.

* `bootable` - Whether the disk is bootable.

* `description` - The disk description.

* `volume_type` - The disk type.

* `multiattach` - Whether the disk is shareable.

* `size` - The disk size, in GB.

* `status` - The disk status.

* `created_at` - The time when the disk was created.

* `updated_at` - The time when the disk was updated.

* `wwn` - The unique identifier used when attaching the disk.

* `metadata` - Metadata key/value pairs to associate with the disk.

The `attachments` block supports:

* `id` - The ID of the attached resource in UUID format.

* `attached_at` - The time when the disk was attached.

* `attached_mode` - The ID of the attachment information.

* `device_name` - The device name to which the disk is attached.

* `server_id` - The ID of the server to which the disk is attached.
