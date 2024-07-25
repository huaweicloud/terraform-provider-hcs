---
subcategory: "Elastic Cloud Server (ECS)"
---

# hcs_ecs_compute_interface_attach

Attaches a Network Interface to an Instance.

## Example Usage

### Attach a port (under the specified network) to the ECS instance and generate a random IP address

```hcl
variable "instance_id" {}
variable "network_id" {}

resource "hcs_ecs_compute_interface_attach" "test" {
  instance_id = var.instance_id
  network_id  = var.network_id
}
```

### Attach a port (under the specified network) to the ECS instance and use the custom security groups

```hcl
variable "instance_id" {}
variable "network_id" {}

resource "hcs_ecs_compute_interface_attach" "test" {
  instance_id        = var.instance_id
  network_id         = var.network_id
  fixed_ip           = "192.168.10.199"
}
```

### Attach a custom port to the ECS instance

```hcl
variable "instance_id" {}

variable "port_id" {}

resource "hcs_ecs_compute_interface_attach" "attached" {
  instance_id = var.instance_id
  port_id = var.port_id
  
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the network interface attache resource. If
  omitted, the provider-level region will be used. Changing this creates a new network interface attache resource.

* `instance_id` - (Required, String, ForceNew) The ID of the Instance to attach the Port or Network to.

* `port_id` - (Optional, String, ForceNew) The ID of the Port to attach to an Instance.
  This option and `network_id` are mutually exclusive.

* `network_id` - (Optional, String, ForceNew) The ID of the Network to attach to an Instance. A port will be created
  automatically.
  This option and `port_id` are mutually exclusive.

* `fixed_ip` - (Optional, String, ForceNew) An IP address to assosciate with the port.

  ->This option cannot be used with port_id. You must specify a network_id. The IP address must lie in a range on
  the supplied network.

* `source_dest_check` - (Optional, Bool, ForceNew) Specifies whether the ECS processes only traffic that is destined 
  specifically for it. This function is enabled by default but should be disabled if the ECS functions as a SNAT server or has a
  virtual IP address bound to it.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of ECS instance ID and port ID separated by a slash.
* `mac` - The MAC address of the NIC.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `delete` - Default is 10 minute.

## Import

Interface Attachments can be imported using the Instance ID and Port ID separated by a slash, e.g.

```shell
$ terraform import hcs_ecs_compute_interface_attach.ai_1 89c60255-9bd6-460c-822a-e2b959ede9d2/45670584-225f-46c3-b33e-6707b589b666
```
