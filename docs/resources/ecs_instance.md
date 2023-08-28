---

# hcs_ecs_compute_instance

Manages a ECS instance resource within HCS.

## Example Usage

### Basic Instance

```hcl
data "hcs_availability_zones" "test" {
  provider = huaweicloudstack
}

resource "hcs_vpc" "vpc" {
  provider = huaweicloudstack
  name = "tf_vpc_test"
  cidr = "192.168.0.0/16"
}

resource "hcs_vpc_subnet" "subnet" {
  provider = huaweicloudstack
  name       = "subnet_1"
  cidr       = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id     = hcs_vpc.vpc.id
}

resource "hcs_networking_secgroup_rule" "test" {
  provider = huaweicloudstack
  security_group_id       = hcs_networking_secgroup.secgroup.id
  direction               = "ingress"
  action                  = "allow"
  ethertype               = "IPv4"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "hcs_networking_secgroup" "secgroup" {
  provider = huaweicloudstack
  name        = "secgroup_1"
  description = "My security group"
}

data "hcs_ims_images" "centos" {
  provider = huaweicloudstack
  name       = "222"
  visibility = "public"
}

data "hcs_ecs_compute_flavors" "flavors" {
  provider = huaweicloudstack
  availability_zone = data.hcs_availability_zones.test.names[0]
  cpu_core_count    = 2
  memory_size       = 4
}

resource "hcs_ecs_compute_instance" "instance" {
  provider = huaweicloudstack
  name               = "tf_ecs-test2"
  image_id           = data.hcs_ims_images.centos.images[0].id
  flavor_id = data.hcs_ecs_compute_flavors.flavors.ids[0]
  security_groups = [hcs_networking_secgroup.secgroup.name]
  availability_zone  = data.hcs_availability_zones.test.names[0]

  network {
    uuid = hcs_vpc_subnet.subnet.id
  }

  block_device_mapping_v2 {
    source_type  = "image"
    destination_type = "volume"
    uuid = data.hcs_ims_images.centos.images[0].id
    volume_type = "business_type_01"
    volume_size = 20
  }
}
```

### Instance With Attached Volume

data "hcs_availability_zones" "test" {
  provider = huaweicloudstack
}

resource "hcs_vpc" "vpc" {
  provider = huaweicloudstack
  name = "tf_vpc"
  cidr = "192.168.0.0/16"
}

resource "hcs_vpc_subnet" "subnet" {
  provider = huaweicloudstack
  name       = "subnet_10010"
  cidr       = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id     = hcs_vpc.vpc.id
}

resource "hcs_networking_secgroup_rule" "test" {
  provider = huaweicloudstack
  security_group_id       = hcs_networking_secgroup.secgroup.id
  direction               = "ingress"
  action                  = "allow"
  ethertype               = "IPv4"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "hcs_networking_secgroup" "secgroup" {
  provider = huaweicloudstack
  name        = "secgroup_10010"
  description = "My security group"
}

data "hcs_ims_images" "centos" {
  provider = huaweicloudstack
  name       = "mini_image"
  visibility = "public"
}

data "hcs_ecs_compute_flavors" "flavors" {
  provider = huaweicloudstack
  availability_zone = data.hcs_availability_zones.test.names[0]
  cpu_core_count    = 2
  memory_size       = 4
}

resource "hcs_ecs_compute_instance" "instance" {
  provider = huaweicloudstack
  name               = "tf_ecs-test2"
  image_id           = data.hcs_ims_images.centos.images[0].id
  flavor_id = data.hcs_ecs_compute_flavors.flavors.ids[0]
  security_groups = [hcs_networking_secgroup.secgroup.name]
  availability_zone  = data.hcs_availability_zones.test.names[0]

  network {
    uuid = hcs_vpc_subnet.subnet.id
  }
  block_device_mapping_v2 {
    source_type  = "image"
    destination_type = "volume"
    uuid = data.hcs_ims_images.centos.images[0].id
    volume_type = "business_type_01"
    volume_size = 20
    }
}

resource "hcs_ecs_compute_volume_attach" "attached" {
  provider = huaweicloudstack
  instance_id = hcs_ecs_compute_instance.instance.id
  volume_id   = "0f2ee145-adfe-4270-b126-ad1e49a6f775"
  device = "/dev/vdb"
}
```

### Instance With Multiple Networks

```hcl
data "hcs_availability_zones" "test" {
  provider = huaweicloudstack
}

resource "hcs_vpc" "vpc" {
  provider = huaweicloudstack
  name = "tf_vpc_10086"
  cidr = "192.168.0.0/16"
}

resource "hcs_vpc_subnet" "subnet" {
  provider = huaweicloudstack
  name       = "subnet_10086"
  cidr       = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id     = hcs_vpc.vpc.id
}

resource "hcs_networking_secgroup_rule" "test" {
  provider = huaweicloudstack
  security_group_id       = hcs_networking_secgroup.secgroup.id
  direction               = "ingress"
  action                  = "allow"
  ethertype               = "IPv4"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "hcs_networking_secgroup" "secgroup" {
  provider = huaweicloudstack
  name        = "secgroup_10086"
  description = "My security group"
}

data "hcs_ims_images" "centos" {
  provider = huaweicloudstack
  name       = "mini_image"
  visibility = "public"
}

data "hcs_ecs_compute_flavors" "flavors" {
  provider = huaweicloudstack
  availability_zone = data.hcs_availability_zones.test.names[0]
  cpu_core_count    = 2
  memory_size       = 4
}

resource "hcs_ecs_compute_instance" "instance" {
  provider = huaweicloudstack
  name               = "tf_ecs-test3"
  image_id           = data.hcs_ims_images.centos.images[0].id
  flavor_id = data.hcs_ecs_compute_flavors.flavors.ids[0]
  security_groups = [hcs_networking_secgroup.secgroup.name]
  availability_zone  = data.hcs_availability_zones.test.names[0]

  network {
    uuid = hcs_vpc_subnet.subnet.id
  }
  
  network {
    uuid = hcs_vpc_subnet.subnet.id
  }
  
  block_device_mapping_v2 {
    source_type  = "image"
    destination_type = "volume"
    uuid = data.hcs_ims_images.centos.images[0].id
    volume_type = "business_type_01"
    volume_size = 20
  }
}
```

### Instance with User Data (cloud-init)

```hcl
data "hcs_availability_zones" "test" {
  provider = huaweicloudstack
}

resource "hcs_vpc" "vpc" {
  provider = huaweicloudstack
  name = "tf_vpc_10086"
  cidr = "192.168.0.0/16"
}

resource "hcs_vpc_subnet" "subnet" {
  provider = huaweicloudstack
  name       = "subnet_10086"
  cidr       = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id     = hcs_vpc.vpc.id
}

resource "hcs_networking_secgroup_rule" "test" {
  provider = huaweicloudstack
  security_group_id       = hcs_networking_secgroup.secgroup.id
  direction               = "ingress"
  action                  = "allow"
  ethertype               = "IPv4"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "hcs_networking_secgroup" "secgroup" {
  provider = huaweicloudstack
  name        = "secgroup_10086"
  description = "My security group"
}

data "hcs_ims_images" "centos" {
  provider = huaweicloudstack
  name       = "ecs_cloudinit_image"
  visibility = "public"
}

data "hcs_ecs_compute_flavors" "flavors" {
  provider = huaweicloudstack
  availability_zone = data.hcs_availability_zones.test.names[0]
  cpu_core_count    = 2
  memory_size       = 4
}

resource "hcs_ecs_compute_instance" "instance" {
  provider = huaweicloudstack
  name               = "tf_ecs-test3"
  image_id           = data.hcs_ims_images.centos.images[0].id
  flavor_id = data.hcs_ecs_compute_flavors.flavors.ids[0]
  security_groups = [hcs_networking_secgroup.secgroup.name]
  availability_zone  = data.hcs_availability_zones.test.names[0]
  user_data       = "xxxxxxxxxxxxxxxxxxxxxxx"

network {
  uuid = hcs_vpc_subnet.subnet.id
}

block_device_mapping_v2 {
  source_type  = "image"
  destination_type = "volume"
  uuid = data.hcs_ims_images.centos.images[0].id
  volume_type = "business_type_01"
  volume_size = 20
}
}
```

`user_data` can come from a variety of sources: inline, read in from the `file`
function, or the `template_cloudinit_config` resource.

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) A unique name for the instance.

* `image_id` - (Required, String, ForceNew) The ID of the desired image for the server. Changing this creates a new
  server.

* `flavor_id` - (Required, String) The ID of the desired flavor for the server. Changing this resizes the existing
  server.

* `user_data` - (Optional, String, ForceNew) The user data to provide when launching the instance.The string length 
   must be less than 65535 and must be encrypted using Base64. Changing this creates a new server.

* `admin_pass` - (Optional, String, ForceNew) The administrative password to assign to the server. Changing this creates
  a new server.

* `key_name` - (Optional, String, ForceNew) The name of a key pair to put on the server. The key pair must already be
  created and associated with the tenant's account. Changing this creates a new server.

* `vpc_id` - (Required, String, ForceNew) The ID of the desired VPC for the server. Changing this creates a new server.

* `network` - (Optional, List, ForceNew) An array of one or more networks to attach to the instance. The network object
  structure is documented below. Changing this creates a new server.

* `system_disk_type` - (Optional, String, ForceNew) The system disk type of the server.
  Changing this creates a new server.

* `system_disk_size` - (Optional, Int, ForceNew) The system disk size in GB, The value range is 1 to 1024. Changing this
  creates a new server.  

* `block_device_mapping_v2` - (Optional, List, ForceNew) An array of one or more data disks to attach to the instance.
  The data_disks object structure is documented below. Changing this creates a new server.

* `security_groups` - (Optional, String) An array of one or more security group names to associate with the server.
  Changing this results in adding/removing security groups from the existing server.

* `availability_zone` - (Required, String, ForceNew) The availability zone in which to create the server.
  Changing this creates a new server.
  
* `delete_disks_on_termination` - (Optional, Bool) Delete the data disks upon termination of the instance. Defaults to
  false. Changing this creates a new server.

* `enterprise_project_id` - (Optional, String) The enterprise project id. Changing this creates a new server.

* `tags` - (Optional, Map) Tags key/value pairs to associate with the instance.

* `op_svc_userid` - (Optional, String, ForceNew) User ID, required when using key_name. Changing this creates a new
  server.

The `network` block supports:

* `uuid` - (Optional, String, ForceNew) The network UUID to attach to the server. Changing this creates a new
  server.

* `port` - (Optional, String, ForceNew) Specifies the IP address of the. Among the three network parameters (port, UUID,
  and fixed_ip), port has the highest priority. The UUID must be specified when fixed_ip is specified.. Changing this
  creates a new server.

The `block_device_mapping_v2` block supports:

* `boot_index` - (Optional, Int, ForceNew)Boot flag. The value 0 indicates the boot disk, and the value - 1 indicates 
  the non-boot disk. Note: When the source types of all volume devices are volume, one value of boot_index is 0.
  Changing this creates a new server.
  
* `destination_type` - (Optional, String) Indicates the current type of the volume device. Currently, only the 
  volume type is supported. Changing this creates a new server.

* `source_type` - (Required, String) SSource type of a volume device. Currently, only the volume, image, and 
   snapshot types are supported. If a volume is used to create an ECS, set source_type to volume. If you use an image to
   create an ECS, set source_type to image. To create an ECS using a snapshot, set source_type to snapshot. Note: If the
   source type of a volume device is snapshot and boot_index is 0, the EVS disk corresponding to the snapshot must be a 
   system disk. Changing this creates a new server.

* `uuid` - (Required, String, ForceNew) Specifies the UUID of the volume or snapshot. If source_type is image, the value 
   is the UUID of the image.. Changing this creates a new server.

* `volume_size` - (Optional, Int, ForceNew) Specifies the volume size. The value is an integer. This parameter is 
   mandatory when source_type is set to image and destination_type is set to volume. The unit is GB.
   Changing this creates a new server.

* `device_name` - (Optional, String, ForceNew) Specifies the volume device name. The value is a string of 0 to 255 
   characters and must comply with the regular expression (^/dev/x{0, 1}[a-z]{0, 1}d{0, 1})([a-z]+)[0-9]*$ . Example: 
   /dev/vda; User-specified device_name.The configuration does not take effect. The system generates a device_name by 
   default.

* `delete_on_termination` - (Optional, Bool, ForceNew) Specifies whether to delete the volume when deleting an ECS. 
   The default value is false.

* `volume_type` - (Optional, String, ForceNew) Specifies the volume type. This parameter is used when source_type is 
   set to image and destination_type is set to volume.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minute.
* `update` - Default is 30 minute.
* `delete` - Default is 30 minute.

## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import hcs_ecs_compute_instance.instance_1 d90ce693-5ccf-4136-a0ed-152ce412b6b9
