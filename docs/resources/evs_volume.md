---
subcategory: "Elastic Volume Service (EVS)"
---

# hcs_evs_volume

Manages a volume resource within HuaweiCloudStack.

## Example Usage

```hcl
resource "hcs_evs_volume" "volume_1" {
  availability_zone      = "az"
  name                   = "volume_1"
  description            = "first test volume"
  size                   = 3
  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the disk. If omitted, the `region` argument of
  the provider is used. Changing this creates a new disk.

* `size` - (Required, Int) The size of the disk to create (in gigabytes).

* `availability_zone` - (Optional, String, ForceNew) The availability zone for the disk. Changing this creates a new
  disk.

* `description` - (Optional, String) A description of the disk. Changing this updates the disk's description.

* `image_id` - (Optional, String, ForceNew) The image ID from which to create the disk. Changing this creates a new
  disk.

* `metadata` - (Optional, Map) Metadata key/value pairs to associate with the disk. Changing this updates the existing
  disk metadata.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID which the desired the disk belongs to.

* `name` - (Optional, String) A unique name for the disk. Changing this updates the disk's name.

* `snapshot_id` - (Optional, String, ForceNew) The snapshot ID from which to create the disk. Changing this creates a
  new disk.

* `source_vol_id` - (Optional, String, ForceNew) The disk ID from which to create the disk. Changing this creates a
  new disk.

* `volume_type` - (Optional, String, ForceNew) The type of disk to create. Changing this creates a new disk.

* `multiattach` - (Optional, Bool, ForceNew) Specifies whether the disk is shareable. The default value is false. 
  Changing this creates a new disk.

* `tags` - (Optional, Map) The key/value pairs to associate with the disk.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `attachments` - If a disk is attached to an instance, this attribute will display the Attachment ID, Instance ID, and
  the Device as the Instance sees it. The [object](#attachments_struct) structure is documented below.

* `wwn` - The unique identifier used for mounting the disk.

* `status` - The status of disk.

* `bootable` - Whether the disk is bootable.

* `created_at` - The time when the disk was created.

* `updated_at` - The time when the disk was updated.

<a name="attachments_struct"></a>
The `attachments` block supports:

* `id` - The ID of the attached resource in UUID format.

* `attached_at` - The time when the disk was attached.

* `attached_mode` - The ID of the attachment information.

* `device_name` - The device name to which the disk is attached.

* `server_id` - The ID of the server to which the disk is attached.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `update` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

Volumes can be imported using the `id`, e.g.

```
$ terraform import hcs_evs_volume.volume_1 ea257959-eeb1-4c10-8d33-26f0409a755d
```