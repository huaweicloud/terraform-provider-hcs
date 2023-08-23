---
subcategory: "Image Management Service (IMS)"
---

# hcs_ims_image_share_accepter

Manages an IMS image share accepter resource within HuaweiCloudStack.

## Example Usage

```hcl
variable "image_id" {}

resource "hcs_ims_image_share_accepter" "test" {
  image_id = var.image_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `image_id` - (Required, String, ForceNew) Specifies the ID of the image.

  Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
