---
subcategory: "Elastic Cloud Server (ECS)"
---

# hcs_ecs_compute_volume_attach

Attaches a volume to an ECS Instance.

## Example Usage

### Basic attachment of a single volume to a single instance

```hcl
variable "instance_id" {}

variable "volume_id" {}

variable "device" {}

resource "hcs_ecs_compute_volume_attach" "attached" {
  instance_id = var.instance_id
  volume_id   = var.volume_id
  device = var.device
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the volume resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the Instance to attach the Volume to.

* `volume_id` - (Required, String, ForceNew) Specifies the ID of the Volume to attach to an Instance.

* `device` - (Optional, String, ForceNew) Specifies the device of the volume attachment (ex: `/dev/vdc`).

  -> Being able to specify a device is dependent upon the hypervisor in use. There is a chance that the device
  specified in Terraform will not be the same device the hypervisor chose. If this happens, Terraform will wish to
  update the device upon subsequent applying which will cause the volume to be detached and reattached indefinitely.
  Please use with caution.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `pci_address` - PCI address of the block device.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

Volume Attachments can be imported using the Instance ID and Volume ID separated by a slash, e.g.

```shell
$ terraform import hcs_ecs_compute_volume_attach.va_1 89c60255-9bd6-460c-822a-e2b959ede9d2/45670584-225f-46c3-b33e-6707b589b666
```
