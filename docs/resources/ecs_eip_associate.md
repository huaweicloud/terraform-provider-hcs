---
subcategory: "Elastic Cloud Server (ECS)"
---

# hcs_ecs_compute_eip_associate

* Associates the **IPv4** address of an ECS instance to a specified EIP.
* Associates the **IPv6** address of an ECS instance to a specified **Shared** Bandwidth.

## Example Usage

### Automatically detect the correct network

```hcl
var "bandwidth_name" {}

resource "hcs_ecs_compute_instance" "myinstance" {
  name = "instance"
  ...

  network {
    uuid = "55534eaa-533a-419d-9b40-ec427ea7195a"
  }
}

resource "hcs_vpc_eip" "myeip" {
  publicip {
    type = "external_network_eip"
  }

  bandwidth {
    share_type  = "PER"
    name        = var.bandwidth_name
    size        = 10
  }
}

resource "hcs_ecs_compute_eip_associate" "associated" {
  public_ip   = hcs_vpc_eip.myeip.address
  instance_id = hcs_ecs_compute_instance.myinstance.id
}
```

### Explicitly set the network to attach to

```hcl
resource "hcs_ecs_compute_instance" "myinstance" {
  name = "instance"
  ...

  network {
    uuid = "55534eaa-533a-419d-9b40-ec427ea7195a"
  }

  network {
    uuid = "3c4a0d74-24b9-46cf-9d7f-8b7a4dc2f65c"
  }
}

resource "hcs_vpc_eip" "myeip" {
  publicip {
    type = "eip"
  }
  bandwidth {
    name        = "test"
    size        = 5
    share_type  = "PER"
  }
}

resource "hcs_ecs_compute_eip_associate" "associated" {
  public_ip   = hcs_vpc_eip.myeip.address
  instance_id = hcs_ecs_compute_instance.myinstance.id
  fixed_ip    = hcs_ecs_compute_instance.myinstance.network.1.fixed_ip_v4
}
```

### Associate the IPv6 address to a specified Shared Bandwidth

```hcl
variable "subnet_id" {}
variable "bandwidth_id" {}

resource "hcs_ecs_compute_instance" "myinstance" {
  name      = "instance"
  flavor_id = "c6.large.2"
  ...

  network {
    uuid        = var.subnet_id
    ipv6_enable = true
  }
}

resource "hcs_ecs_compute_eip_associate" "associated" {
  bandwidth_id = var.bandwidth_id
  instance_id  = hcs_ecs_compute_instance.myinstance.id
  fixed_ip     = hcs_ecs_compute_instance.myinstance.network.0.fixed_ip_v6
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the associated resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of ECS instance to associated with.
  Changing this creates a new resource.

* `public_ip` - (Optional, String, ForceNew) Specifies the EIP address to associate. It's **mandatory**
  when you want to associate the ECS instance with an EIP. Changing this creates a new resource.

* `bandwidth_id` - (Optional, String, ForceNew) Specifies the **shared** bandwidth ID to associate.
  It's **mandatory** when you want to associate the ECS instance with a specified shared bandwidth.
  Changing this creates a new resource.

* `fixed_ip` - (Optional, String, ForceNew) Specifies the private IP address to direct traffic to. It's **mandatory**
  and must be a valid IPv6 address when you want to associate the ECS instance with a specified shared bandwidth.
  Changing this creates a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of `<eip address or bandwidth_id>/<instance_id>/<fixed_ip>`.
* `port_id` - The port ID of the ECS instance that associated with.

## Import

This resource can be imported by specifying all three arguments, separated by a forward slash:

```shell
$ terraform import hcs_ecs_compute_eip_associate.bind <eip address or bandwidth_id>/<instance_id>/<fixed_ip>
```
