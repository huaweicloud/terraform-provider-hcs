---
subcategory: "Image Management Service (IMS)"
---

# hcs_ims_images

Use this data source to get available HuaweiCloudStack IMS images.

## Example Usage

```hcl
data "hcs_ims_images" "centos" {
  name       = "CentOS_7.4_64bit"
  visibility = "public"
}

data "hcs_ims_images" "centos-1" {
  architecture = "x86_64"
  os_version   = "CentOS 7.4 64bit"
  visibility   = "public"
}

data "hcs_ims_images" "centos-2" {
  architecture = "x86_64"
  name_regex   = "^CentOS_7.4"
  visibility   = "public"
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the images. If omitted, the provider-level region will be
  used.

* `name` - (Optional, String) The name of the image. Cannot be used simultaneously with `name_regex`.

* `name_regex` - (Optional, String) The regular expression of the name of the image.
  Cannot be used simultaneously with `name`.

* `visibility` - (Optional, String) The visibility of the image. Must be one of
  **public**, **private**, or **shared**.

* `architecture` - (Optional, String) Specifies the image architecture type. The value can be **x86_64** and **aarch64**.

* `os` - (Optional, String) Specifies the image OS type. The value can be **Windows**, **Ubuntu**,
  **RedHat**, **SUSE**, **CentOS**, **Debian**, **OpenSUSE**, **Oracle Linux**, **Fedora**, **Other**,
  **CoreOS**, or **EulerOS**.

* `os_version` - (Optional, String) Specifies the OS version. For example, *CentOS 7.4 64bit* or *Ubuntu 18.04 server
  64bit*.

* `image_type` - (Optional, String) Specifies the environment where the image is used. For a BMS image, the value is *
  *Ironic**.

* `owner` - (Optional, String) The owner (UUID) of the image.

* `tag` - (Optional, String) Search for images with a specific tag in "Key=Value" format.

* `sort_direction` - (Optional, String) Order the results in either `asc` or `desc`.

* `sort_key` - (Optional, String) Sort images based on a certain key. Must be one of
  "name", "container_format", "disk_format", "status", "id" or "size". Defaults to `name`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `images` - Indicates the images' information. Structure is documented below.

The `images` block contains:

* `id` - The ID of the image

* `name` - The name of the image.

* `visibility` - The visibility of the image.

* `checksum` - The checksum of the data associated with the image.

* `container_format` - The format of the image's container.

* `disk_format` - The format of the image's disk.

* `min_disk_gb` - The minimum amount of disk space required to use the image.

* `min_ram_mb` - The minimum amount of ram required to use the image.

* `owner` - The owner (UUID) of the image.

* `protected` - Whether the image is protected.

* `image_type` - The environment where the image is used. For a BMS image, the value is **Ironic**.

* `os` - Specifies the image OS type.

* `os_version` - The OS version.

* `status` - The status of the image.

* `created_at` - The date when the image was created.

* `updated_at` - The date when the image was last updated.

* `size_bytes` - The size of the image (in bytes).
