---
subcategory: "Elastic Cloud Server (ECS)"
---

# hcs_ecs_compute_instance

Manages a ECS VM instance resource within HCS.

## Example Usage

### Basic Instance

```hcl
data "hcs_availability_zones" "test" {
}

data "hcs_ecs_compute_flavors" "flavors" {
  availability_zone = data.hcs_availability_zones.test.names[0]
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcs_vpc_subnets" "test" {
  name = "subnet-32a8"
}

data "hcs_ims_images" "test" {
  name       = "mini_image"
}

data "hcs_networking_secgroups" "test" {
  name = "default"
}

resource "hcs_ecs_compute_instance" "test" {
  name                = "tf_ecs-test1"
  description         = "terraform test"
  image_id            = data.hcs_ims_images.test.images[0].id
  flavor_id           = data.hcs_ecs_compute_flavors.flavors.ids[0]
  security_group_ids  = [data.hcs_networking_secgroups.test.security_groups[0].id]
  availability_zone = data.hcs_availability_zones.test.names[0]

  network {
    uuid              = data.hcs_vpc_subnets.test.subnets[0].id
    source_dest_check = false
  }

  system_disk_type = "business_type_01"
  system_disk_size = 10

  data_disks {
    type = "business_type_01"
    size = "10"
  }
  delete_disks_on_termination = true
  delete_eip_on_termination = true
}
```

### Instance With Associated Eip

```hcl
data "hcs_availability_zones" "test" {
}

data "hcs_ecs_compute_flavors" "flavors" {
  availability_zone = data.hcs_availability_zones.test.names[0]
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcs_vpc_subnets" "test" {
  name = "subnet-32a8"
}

data "hcs_ims_images" "test" {
  name       = "mini_image"
}

data "hcs_networking_secgroups" "test" {
  name = "default"
}

resource "hcs_ecs_compute_instance" "test" {
  name                = "tf_ecs-test2"
  description         = "terraform test"
  image_id            = data.hcs_ims_images.test.images[0].id
  flavor_id           = data.hcs_ecs_compute_flavors.flavors.ids[0]
  security_group_ids  = [data.hcs_networking_secgroups.test.security_groups[0].id]
  availability_zone = data.hcs_availability_zones.test.names[0]

  network {
    uuid              = data.hcs_vpc_subnets.test.subnets[0].id
    source_dest_check = false
  }

  system_disk_type = "business_type_01"
  system_disk_size = 10

  data_disks {
    type = "business_type_01"
    size = "10"
  }
  delete_disks_on_termination = true
  delete_eip_on_termination = true
}

resource "hcs_vpc_eip" "test" {
  publicip {
    type = "eip"
  }
  bandwidth {
    name        = "test"
    size        = 100
    share_type  = "PER"
  }
}

resource "hcs_ecs_compute_eip_associate" "associated" {
  public_ip   = hcs_vpc_eip.test.address
  instance_id = hcs_ecs_compute_instance.test.id
}
```

### Instance With Attached Volume

```hcl
data "hcs_availability_zones" "test" {
}

data "hcs_ecs_compute_flavors" "flavors" {
  availability_zone = data.hcs_availability_zones.test.names[0]
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcs_vpc_subnets" "test" {
  name = "subnet-32a8"
}

data "hcs_ims_images" "test" {
  name       = "mini_image"
}

data "hcs_networking_secgroups" "test" {
  name = "default"
}

resource "hcs_ecs_compute_instance" "ecs-volume-attached" {
  name                = "tf_ecs-test3"
  description         = "terraform test"
  image_id            = data.hcs_ims_images.test.images[0].id
  flavor_id           = data.hcs_ecs_compute_flavors.flavors.ids[0]
  security_group_ids  = [data.hcs_networking_secgroups.test.security_groups[0].id]
  availability_zone = data.hcs_availability_zones.test.names[0]

  network {
    uuid              = data.hcs_vpc_subnets.test.subnets[0].id
    source_dest_check = false
  }

  system_disk_type = "business_type_01"
  system_disk_size = 10

  data_disks {
    type = "business_type_01"
    size = "10"
  }
  delete_disks_on_termination = true
  delete_eip_on_termination = true
}

data "hcs_evs_volumes" "test" {
  status = "available"
}

resource "hcs_evs_volume" "test" {
  availability_zone      = data.hcs_availability_zones.test.names[0]
  name                   = "ecs-volume-attached-data"
  description            = "ecs data volume"
  size                   = 10
}

resource "hcs_ecs_compute_volume_attach" "attached" {
  instance_id = hcs_ecs_compute_instance.ecs-volume-attached.id
  volume_id   = hcs_evs_volume.test.id
  device = "/dev/vdb"
}
```

### Instance With Multiple Data Disks

It's possible to specify multiple `data_disks` entries to create an instance with multiple data disks, but we can't
ensure the volume attached order. So it's recommended to use `Instance With Attached Volume` above.

```hcl
data "hcs_availability_zones" "test" {
}

data "hcs_ecs_compute_flavors" "flavors" {
  availability_zone = data.hcs_availability_zones.test.names[0]
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcs_vpc_subnets" "test" {
  name = "subnet-32a8"
}

data "hcs_ims_images" "test" {
  name       = "mini_image"
}

data "hcs_networking_secgroups" "test" {
  name = "default"
}

resource "hcs_ecs_compute_instance" "multi-disk" {
  name                = "ecs-multi-disk"
  description         = "terraform test"
  image_id            = data.hcs_ims_images.test.images[0].id
  flavor_id           = data.hcs_ecs_compute_flavors.flavors.ids[0]
  security_group_ids  = [data.hcs_networking_secgroups.test.security_groups[0].id]
  availability_zone = data.hcs_availability_zones.test.names[0]

  system_disk_type = "business_type_01"
  system_disk_size = 10

  data_disks {
    type = "business_type_01"
    size = "10"
  }
  data_disks {
    type = "business_type_01"
    size = "20"
  }

  delete_disks_on_termination = true
  delete_eip_on_termination = true

  network {
    uuid = data.hcs_vpc_subnets.test.subnets[0].id
    source_dest_check = false
  }
}
```

### Instance With Multiple Networks

```hcl
data "hcs_availability_zones" "test" {
}

data "hcs_ecs_compute_flavors" "flavors" {
  availability_zone = data.hcs_availability_zones.test.names[0]
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcs_vpc_subnets" "test" {
  name = "subnet-32a8"
}

data "hcs_ims_images" "test" {
  name       = "mini_image"
}

data "hcs_networking_secgroups" "test" {
  name = "default"
}

resource "hcs_ecs_compute_instance" "instance" {
  name               = "ecs-muti-net"
  image_id           = data.hcs_ims_images.centos.images[0].id
  flavor_id = data.hcs_ecs_compute_flavors.flavors.ids[0]
  security_groups = [hcs_networking_secgroup.secgroup.name]
  availability_zone  = data.hcs_availability_zones.test.names[0]

  network {
    uuid = data.hcs_vpc_subnets.test.subnets[0].id
  }

  network {
    uuid = hcs_vpc_subnet.subnet.id
  }

  delete_disks_on_termination = true
  delete_eip_on_termination = true

  system_disk_type = "business_type_01"
  system_disk_size = 10
}
```

### Instance with User Data (cloud-init)

```hcl
data "hcs_availability_zones" "test" {
}

data "hcs_ecs_compute_flavors" "flavors" {
  availability_zone = data.hcs_availability_zones.test.names[0]
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcs_vpc_subnets" "test" {
  name = "subnet-32a8"
}

data "hcs_ims_images" "test" {
  name       = "mini_image"
}

data "hcs_networking_secgroups" "test" {
  name = "default"
}

resource "hcs_ecs_compute_instance" "ecs-userdata" {
  name                = "ecs-userdata"
  description         = "terraform test"
  image_id            = data.hcs_ims_images.test.images[0].id
  flavor_id           = data.hcs_ecs_compute_flavors.flavors.ids[0]
  security_group_ids  = [data.hcs_networking_secgroups.test.security_groups[0].id]
  availability_zone = data.hcs_availability_zones.test.names[0]
  user_data       = "xxxxxxxxxxxxxxxxxxxxxxx"
  
  network {
    uuid              = data.hcs_vpc_subnets.test.subnets[0].id
    source_dest_check = false
  }

  system_disk_type = "business_type_01"
  system_disk_size = 10

  data_disks {
    type = "business_type_01"
    size = "10"
  }
  delete_disks_on_termination = true
  delete_eip_on_termination = true
}
```
### Instance with Encrypt Data Volumes
```
data "hcs_availability_zones" "test" {
}

data "hcs_ecs_compute_flavors" "flavors" {
  availability_zone = data.hcs_availability_zones.test.names[0]
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcs_vpc_subnets" "test" {
  name = "subnet-32a8"
}

data "hcs_ims_images" "test" {
  name       = "mini_image"
}

data "hcs_networking_secgroups" "test" {
  name = "default"
}

resource "hcs_ecs_compute_instance" "ecs-userdata" {
  name                = "ecs-userdata"
  description         = "terraform test"
  image_id            = data.hcs_ims_images.test.images[0].id
  flavor_id           = data.hcs_ecs_compute_flavors.flavors.ids[0]
  security_group_ids  = [data.hcs_networking_secgroups.test.security_groups[0].id]
  availability_zone = data.hcs_availability_zones.test.names[0]
  user_data       = "xxxxxxxxxxxxxxxxxxxxxxx"
  
  network {
    uuid              = data.hcs_vpc_subnets.test.subnets[0].id
    source_dest_check = false
  }

  system_disk_type = "business_type_01"
  system_disk_size = 10

  data_disks {
    kms_key_id = "ce488d6a-6090-4f7f-a95b-4faf3ce0bad0"
    encrypt_cipher = "AES256-XTS"
    type = "business_type_01"
    size = "10"
  }
  delete_disks_on_termination = true
  delete_eip_on_termination = true
}
```

### Instance with Encrypt Sys Volumes
```
data "hcs_availability_zones" "test" {
}

data "hcs_ecs_compute_flavors" "flavors" {
  availability_zone = data.hcs_availability_zones.test.names[0]
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcs_vpc_subnets" "test" {
  name = "subnet-32a8"
}

data "hcs_ims_images" "test" {
  name       = "mini_image"
}

data "hcs_networking_secgroups" "test" {
  name = "default"
}

resource "hcs_ecs_compute_instance" "ecs-userdata" {
  name                = "ecs-userdata"
  description         = "terraform test"
  image_id            = data.hcs_ims_images.test.images[0].id
  flavor_id           = data.hcs_ecs_compute_flavors.flavors.ids[0]
  security_group_ids  = [data.hcs_networking_secgroups.test.security_groups[0].id]
  availability_zone = data.hcs_availability_zones.test.names[0]
  user_data       = "xxxxxxxxxxxxxxxxxxxxxxx"
  
  network {
    uuid              = data.hcs_vpc_subnets.test.subnets[0].id
    source_dest_check = false
  }

  system_disk_type = "business_type_01"
  system_disk_size = 10
  kms_key_id = "ce488d6a-6090-4f7f-a95b-4faf3ce0bad0"
  encrypt_cipher = "AES256-XTS"
  
  data_disks {
    kms_key_id = "ce488d6a-6090-4f7f-a95b-4faf3ce0bad0"
    encrypt_cipher = "AES256-XTS"
    type = "business_type_01"
    size = "10"
  }
  delete_disks_on_termination = true
  delete_eip_on_termination = true
}
```

### Instance with power off
```
data "hcs_availability_zones" "test" {
}

data "hcs_ecs_compute_flavors" "flavors" {
  availability_zone = data.hcs_availability_zones.test.names[0]
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcs_vpc_subnets" "test" {
  name = "subnet-32a8"
}

data "hcs_ims_images" "test" {
  name       = "mini_image"
}

data "hcs_networking_secgroups" "test" {
  name = "default"
}

resource "hcs_ecs_compute_instance" "ecs-power-off" {
  name                = "ecs-poweroff"
  description         = "poweroff test"
  image_id            = data.hcs_ims_images.test.images[0].id
  flavor_id           = data.hcs_ecs_compute_flavors.flavors.ids[0]
  security_group_ids  = [data.hcs_networking_secgroups.test.security_groups[0].id]
  availability_zone = data.hcs_availability_zones.test.names[0]
  
  network {
    uuid              = data.hcs_vpc_subnets.test.subnets[0].id
    source_dest_check = false
  }

  system_disk_type = "business_type_01"
  system_disk_size = 10
  
  data_disks {
    type = "business_type_01"
    size = "10"
  }
  
  power_action = "OFF"
}
```

### Instance with tags
```
data "hcs_availability_zones" "test" {
}

data "hcs_ecs_compute_flavors" "test" {
  availability_zone = data.hcs_availability_zones.test.names[0]
  cpu_core_count    = 2
  memory_size       = 4
}

data "hcs_vpc_subnets" "test" {
  name = var.subnet_name
}

data "hcs_ims_images" "test" {
  name       = var.image_name
}

data "hcs_networking_secgroups" "test" {
  name = var.secgroup_name
}

resource "hcs_ecs_compute_instance" "test_with_tags" {
  name                = var.ecs_name
  description         = var.ecs_description
  image_id            = data.hcs_ims_images.test.images[0].id
  flavor_id           = data.hcs_ecs_compute_flavors.test.ids[0]
  security_group_ids  = [data.hcs_networking_secgroups.test.security_groups[0].id]
  availability_zone = data.hcs_availability_zones.test.names[0]
  network {
    uuid              = data.hcs_vpc_subnets.test.subnets[0].id
    source_dest_check = false
  }
  system_disk_type = var.disk_type
  system_disk_size = var.system_disk_size

  data_disks {
    type = var.disk_type
    size = var.data_disk_size
  }
  delete_disks_on_termination = true
  delete_eip_on_termination = true
  tags = {
    tag_key = "tag_value"
    ecs_tag_key = "ecs_tag_value"
  }
}
```
## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the instance.
  If omitted, the provider-level region will be used. Changing this creates a new instance.

* `name` - (Required, String) Specifies a unique name for the instance. The name consists of 1 to 64 characters,
  including letters, digits, underscores (_), hyphens (-), and periods (.).

* `flavor_id` - (Required, String) Specifies the flavor ID of the instance to be created.

* `image_id` - (Optional, String, ForceNew) Required if `image_name` is empty. Specifies the image ID of the desired
  image for the instance. Changing this creates a new instance.

* `image_name` - (Optional, String, ForceNew) Required if `image_id` is empty. Specifies the name of the desired image
  for the instance. Changing this creates a new instance.

* `security_group_ids` - (Optional, List) Specifies an array of one or more security group IDs to associate with the
  instance.

* `availability_zone` - (Optional, String, ForceNew) Specifies the availability zone in which to create the instance.

* `network` - (Required, List, ForceNew) Specifies an array of one or more networks to attach to the instance. The
  network object structure is documented below. Changing this creates a new instance.

* `description` - (Optional, String) Specifies the description of the instance. The description consists of 0 to 85
  characters, and can't contain '<' or '>'.

* `admin_pass` - (Optional, String) Specifies the administrative password to assign to the instance, conflict with `key_pair`.

* `key_pair` - (Optional, String) Specifies the SSH keypair name used for logging in to the instance.

* `private_key` - (Optional, String) Specifies the the private key of the keypair in use. This parameter is mandatory
  when replacing or unbinding a keypair and the instance is in **Running** state.

* `system_disk_type` - (Optional, String, ForceNew) Specifies the system disk type of the instance.
  Changing this creates a new instance.

* `system_disk_size` - (Optional, Int) Specifies the system disk size in GB, The value range is 1 to 1024.
  Shrinking the disk is not supported.

* `data_disks` - (Optional, List, ForceNew) Specifies an array of one or more data disks to attach to the instance.
  The data_disks object structure is documented below. Changing this creates a new instance.

* `eip_type` - (Optional, String, ForceNew) Specifies the type of an EIP that will be automatically assigned to the instance.
  Available values are *5_bgp* (dynamic BGP) and *5_sbgp* (static BGP). Changing this creates a new instance.

* `bandwidth` - (Optional, List, ForceNew) Specifies the bandwidth of an EIP that will be automatically assigned to the instance.
  The object structure is documented below. Changing this creates a new instance.

* `eip_id` - (Optional, String, ForceNew) Specifies the ID of an *existing* EIP assigned to the instance.
  This parameter and `eip_type`, `bandwidth` are alternative. Changing this creates a new instance.

* `user_data` - (Optional, String, ForceNew) Specifies the user data to be injected during the instance creation. Text
  and text files can be injected. Changing this creates a new instance.

  -> **NOTE:** If the `user_data` field is specified for a Linux ECS that is created using an image with Cloud-Init
  installed, the `admin_pass` field becomes invalid.

* `scheduler_hints` - (Optional, List) Specifies the scheduler with hints on how the instance should be launched. The
  available hints are described below.

* `delete_disks_on_termination` - (Optional, Bool) Specifies whether to delete the data disks when the instance is terminated.
  Defaults to *false*.

* `delete_eip_on_termination` - (Optional, Bool) Specifies whether the EIP is released when the instance is terminated.
  Defaults to *true*.

* `power_action` - (Optional, String) Specifies the power status of the instance. The value must be one of the following: *ON*, *OFF*, *REBOOT*, *FORCE-OFF* and *FORCE-REBOOT*.

* `enterprise_project_id` - (Optional, String) Specifies a unique id in UUID format of enterprise project.

  -> **NOTE:** Spot price ECSs are suitable for stateless, fault-tolerant instances that are not sensitive to
  interruptions because they can be reclaimed suddenly. When the market price is higher than the maximum price
  you specified, or the inventory is insufficient, your spot ECS will be terminated.
  Do not use a spot ECS for inflexible or long-term workloads.

* `user_id` - (Optional, String, ForceNew) Specifies a user ID, required when using key_pair in prePaid charging mode.
  Changing this creates a new instance.

* `tags` - (Optional, Map) Tags key/value pairs to associate with the instance.

The `network` block supports:

* `uuid` - (Required, String, ForceNew) Specifies the network UUID to attach to the instance.
  Changing this creates a new instance.

* `fixed_ip_v4` - (Optional, String, ForceNew) Specifies a fixed IPv4 address to be used on this network.
  Changing this creates a new instance.

* `ipv6_enable` - (Optional, Bool, ForceNew) Specifies whether the IPv6 function is enabled for the nic.
  Defaults to false. Changing this creates a new instance.

* `source_dest_check` - (Optional, Bool) Specifies whether the ECS processes only traffic that is destined specifically
  for it. This function is enabled by default but should be disabled if the ECS functions as a SNAT server or has a
  virtual IP address bound to it.

* `access_network` - (Optional, Bool) Specifies if this network should be used for provisioning access.
  Accepts true or false. Defaults to false.

The `data_disks` block supports:

* `type` - (Required, String, ForceNew) Specifies the ECS data disk type, which must be one of available disk types,
  contains of *SSD*, *GPSSD* and *SAS*. Changing this creates a new instance.

* `size` - (Required, Int, ForceNew) Specifies the data disk size, in GB. The value ranges form 10 to 32768.
  Changing this creates a new instance.

* `snapshot_id` - (Optional, String, ForceNew) Specifies the snapshot id. Changing this creates a new instance.

* `kms_key_id` - (Optional, String, ForceNew) Specifies the ID of a KMS key. This is used to encrypt the disk.

* `encrypt_cipher` - (Optional, String, ForceNew) Specifies the encrypt cipher of KMS. This value must be set to *AES256-XTS* or *SM4-XTS* when SM series cryptographic algorithms are used. When other cryptographic algorithms are used, this value must be *AES256-XTS*.
  This param must exist if *kms_key_id* exists

The `bandwidth` block supports:

* `share_type` - (Required, String, ForceNew) Specifies the bandwidth sharing type. Changing this creates a new instance.
  Possible values are as follows:
  + **PER**: Dedicated bandwidth
  + **WHOLE**: Shared bandwidth

* `size` - (Optional, Int, ForceNew) Specifies the bandwidth size. The value ranges from 1 to 300 Mbit/s.
  This parameter is mandatory when `share_type` is set to **PER**. Changing this creates a new instance.

* `id` - (Optional, String, ForceNew) Specifies the **shared** bandwidth id. This parameter is mandatory when
  `share_type` is set to **WHOLE**. Changing this creates a new instance.

* `charge_mode` - (Optional, String, ForceNew) Specifies the bandwidth billing mode. The value can only be zero.

The `scheduler_hints` block supports:

* `group` - (Optional, String, ForceNew) Specifies a UUID of a Server Group.
  The instance will be placed into that group. Changing this creates a new instance.

* `tenancy` - (Optional, String, ForceNew) Specifies the tenancy specifies whether the ECS is to be created on a
  Dedicated Host
  (DeH) or in a shared pool. Changing this creates a new instance.

* `deh_id` - (Optional, String, ForceNew) Specifies the ID of DeH.
  This parameter takes effect only when the value of tenancy is dedicated. Changing this creates a new instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID in UUID format.
* `status` - The status of the instance.
* `system_disk_id` - The system disk voume ID.
* `flavor_name` - The flavor name of the instance.
* `security_groups` - An array of one or more security groups to associate with the instance.
* `public_ip` - The EIP address that is associted to the instance.
* `access_ip_v4` - The first detected Fixed IPv4 address or the Floating IP.
* `access_ip_v6` - The first detected Fixed IPv6 address.
* `created_at` - The creation time, in UTC format.
* `updated_at` - The last update time, in UTC format.

* `network` - An array of one or more networks to attach to the instance.
  The [network object](#compute_instance_network_object) structure is documented below.

* `volume_attached` - An array of one or more disks to attach to the instance.
  The [volume attached object](#compute_instance_volume_object) structure is documented below.

<a name="compute_instance_network_object"></a>
The `network` block supports:

* `port` - The port ID corresponding to the IP address on that network.
* `mac` - The MAC address of the NIC on that network.
* `fixed_ip_v4` - The fixed IPv4 address of the instance on this network.
* `fixed_ip_v6` - The Fixed IPv6 address of the instance on that network.

<a name="compute_instance_volume_object"></a>
The `volume_attached` block supports:

* `volume_id` - The volume ID on that attachment.
* `boot_index` - The volume boot index on that attachment.
* `is_sys_volume` - Whether the volume is the system disk.
* `size` - The volume size on that attachment.
* `type` - The volume type on that attachment.
* `pci_address` - The volume pci address on that attachment.

## Import

Instances can be imported by their `id`. For example,

```
terraform import hcs_ecs_compute_instance.my_instance b11b407c-e604-4e8d-8bc4-92398320b847
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `admin_pass`, `user_data`, `data_disks`, `scheduler_hints`,
`delete_disks_on_termination`, `delete_eip_on_termination`, `network/access_network`, `bandwidth`, `eip_type` and arguments for pre-paid and spot price.
It is generally recommended running `terraform plan` after importing an instance.
You can then decide if changes should be applied to the instance, or the resource definition should be updated to
align with the instance. Also you can ignore changes as below.

```
resource "hcs_ecs_compute_instance" "myinstance" {
    ...

  lifecycle {
    ignore_changes = [
      user_data, data_disks,
    ]
  }
}
```

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minute.
* `update` - Default is 30 minute.
* `delete` - Default is 30 minute.