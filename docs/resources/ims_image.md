---
subcategory: "Image Management Service (IMS)"
---

# hcs_ims_image

Manages an Image resource within HuaweiCloudStack IMS.

## Example Usage

### Creating an image from an existing ECS

```hcl
variable "instance_name" {}
variable "image_name" {}

data "hcs_ecs_compute_instance" "test" {
  name = var.instance_name
}

resource "hcs_ims_image" "test" {
  name        = var.image_name
  instance_id = data.hcs_ecs_compute_instance.test.id
  description = "created by Terraform"
}
```

### Creating an image from OBS bucket

```hcl
resource "hcs_ims_image" "ims_test_file" {
  name        = "ims_test_file"
  image_url   = "ims-image:centos70.qcow2"
  min_disk    = 40
  os_version  = "Other(64 bit) 64bit"
  description = "Create an image from the OBS bucket."
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of the image.

* `description` - (Optional, String, ForceNew) A description of the image.

* `instance_id` - (Optional, String, ForceNew) The ID of the ECS that needs to be converted into an image. This
  parameter is mandatory when you create a privete image from an ECS.

* `image_url` - (Optional, String, ForceNew) The URL of the external image file in the OBS bucket. This parameter is
  mandatory when you create a private image from an external file uploaded to an OBS bucket. The format is *OBS bucket
  name:Image file name*.

* `min_ram` - (Optional, Int, ForceNew) The minimum memory of the image in the unit of MB. The default value is 0,
  indicating that the memory is not restricted.

* `min_disk` - (Optional, Int, ForceNew) The minimum size of the system disk in the unit of GB. This parameter is
  mandatory when you create a private image from an external file uploaded to an OBS bucket. The value ranges from 1 GB
  to 1024 GB.

* `os_version` - (Optional, String, ForceNew) The OS version. This parameter is valid when you create a private image
  from an external file uploaded to an OBS bucket.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A unique ID assigned by IMS.

* `visibility` - Whether the image is visible to other tenants.

* `data_origin` - The image resource. The pattern can be 'instance,*instance_id*' or 'file,*image_url*'.

* `disk_format` - The image file format. The value can be `vhd`, `zvhd`, `raw`, `zvhd2`, or `qcow2`.

* `image_size` - The size(bytes) of the image file format.

* `checksum` - The checksum of the data associated with the image.

* `status` - The status of the image.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `delete` - Default is 3 minute.

## Import

Images can be imported using the `id`, e.g.

```sh
terraform import hcs_ims_image.my_image 7886e623-f1b3-473e-b882-67ba1c35887f
```
