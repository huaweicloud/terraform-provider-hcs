---
subcategory: "Elastic Cloud Server (ECS)"
---

# hcs_ecs_compute_volume_attach

Attaches a volume to an ECS Instance.

## Example Usage

### Basic attachment of a single volume to a single instance

```hcl
variable "security_group_id" {}

resource "hcs_evs_volume" "myvol" {
  name              = "volume"
  availability_zone = "cn-north-4a"
  volume_type       = "SAS"
  size              = 10
}

resource "hcs_ecs_compute_instance" "myinstance" {
  name               = "instance"
  image_id           = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id          = "s6.small.1"
  key_pair           = "my_key_pair_name"
  security_group_ids = [var.security_group_id]
  availability_zone  = "cn-north-4a"

  network {
    uuid = "55534eaa-533a-419d-9b40-ec427ea7195a"
  }
}

resource "hcs_ecs_compute_volume_attach" "attached" {
  instance_id = hcs_ecs_compute_instance.myinstance.id
  volume_id   = hcs_evs_volume.myvol.id
}
```

### Attaching multiple volumes to a single instance

```hcl
variable "security_group_id" {}

resource "hcs_evs_volume" "myvol" {
  count             = 2
  name              = "volume_1"
  availability_zone = "cn-north-4a"
  volume_type       = "SAS"
  size              = 10
}

resource "hcs_ecs_compute_instance" "myinstance" {
  name               = "instance"
  image_id           = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id          = "s6.small.1"
  key_pair           = "my_key_pair_name"
  security_group_ids = [var.security_group_id]
  availability_zone  = "cn-north-4a"
}

resource "hcs_ecs_compute_volume_attach" "attachments" {
  count       = 2
  instance_id = hcs_ecs_compute_instance.myinstance.id
  volume_id   = element(hcs_evs_volume.myvol[*].id, count.index)
}

output "volume devices" {
  value = hcs_ecs_compute_volume_attach.attachments[*].device
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the volume resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the Instance to attach the Volume to.

* `volume_id` - (Required, String, ForceNew) Specifies the ID of the Volume to attach to an Instance.

* `device` - (Optional, String) Specifies the device of the volume attachment (ex: `/dev/vdc`).

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
