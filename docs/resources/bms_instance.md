---
subcategory: "Bare Metal Server (BMS)"
---

# hcs_bms_instance

Manages a BMS instance resource within HuaweiCloudStack.

## Example Usage

### Basic Instance

```hcl
variable "instance_name" {}

variable "image_id" {}

variable "flavor_id" {}

variable "user_id" {}

variable "key_pair" {}

data "hcs_availability_zones" "myaz" {}

data "hcs_vpc" "myvpc" {
  name = "vpc-default"
}

data "hcs_vpc_eip" "myeip" {
  public_ip = "192.168.0.123"
}

data "hcs_vpc_subnet" "mynet" {
  name = "subnet-default"
}

data "hcs_networking_secgroup" "mysecgroup" {
  name = "default"
}

resource "hcs_bms_instance" "test" {
  name                  = var.instance_name
  image_id              = var.image_id
  flavor_id             = var.flavor_id
  security_groups       = [data.hcs_networking_secgroup.mysecgroup.id]
  availability_zone     = data.hcs_availability_zones.myaz.names[0]
  vpc_id                = data.hcs_vpc.myvpc.id
  eip_id                = data.hcs_vpc_eip.myeip.id
  key_pair              = var.key_pair

  data_disks {
    type = "SSD"
    size = 100
  }

  nics {
    subnet_id  = data.hcs_vpc_subnet.mynet.id
    ip_address = "192.168.0.123"
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the instance. If omitted, the
  provider-level region will be used. Changing this creates a new instance.

* `name` - (Required, String) Specifies a unique name for the instance. The name consists of 1 to 63 characters,
  including letters, digits, underscores (_), hyphens (-), and periods (.).

* `image_id` - (Required, String, ForceNew) Specifies the image ID of the desired image for the instance. Changing this
  creates a new instance.

* `flavor_id` - (Required, String, ForceNew) Specifies the flavor ID of the desired flavor for the instance. Changing
  this creates a new instance.

* `user_id` - (Optional, String, ForceNew) Specifies the user ID. You can obtain the user ID from My Credential on the
  management console. Changing this creates a new instance.

* `availability_zone` - (Required, String, ForceNew) Specifies the availability zone in which to create the instance.
  for the values. Changing this creates a new instance.

* `vpc_id` - (Required, String, ForceNew) Specifies id of vpc in which to create the instance. Changing this creates a
  new instance.

* `nics` - (Required, List, ForceNew) Specifies an array of one or more networks to attach to the instance. The network
  object structure is documented below. Changing this creates a new instance.

* `admin_pass` - (Optional, String, ForceNew) Specifies the administrative password to assign to the instance. Changing
  this creates a new instance.

* `key_pair` - (Optional, String, ForceNew) Specifies the name of a key pair to put on the instance. The key pair must
  already be created and associated with the tenant's account. Changing this creates a new instance.

* `user_data` - (Optional, String, ForceNew) Specifies the user data to be injected during the instance creation. Text
  and text files can be injected. `user_data` can come from a variety of sources: inline, read in from the
  *file* function. Changing this creates a new instance.

-> **NOTE:** If the `user_data` field is specified for a Linux BMS that is created using an image with Cloud-Init
installed, the `admin_pass` field becomes invalid.

* `security_groups` - (Optional, List, ForceNew) Specifies an array of one or more security group IDs to associate with
  the instance. Changing this creates a new instance.

* `eip_id` - (Optional, String, ForceNew) The ID of the EIP. Changing this creates a new instance.

-> **NOTE:** If the eip_id parameter is configured, you do not need to configure the bandwidth parameters:
`iptype`, `bandwidth_size`, `share_type`.

* `iptype` - (Optional, String, ForceNew) Elastic IP type. Changing this creates a new instance.

* `sharetype` - (Optional, String, ForceNew) Bandwidth sharing type. Changing this creates a new instance. Available
  options are:
    + `PER`: indicates dedicated bandwidth.
    + `WHOLE`: indicates shared bandwidth.

* `bandwidth_size` - (Optional, Int, ForceNew) Bandwidth size. Changing this creates a new instance.

* `data_disks` - (Optional, List, ForceNew) Specifies an array of one or more data disks to attach to the instance. The
  data_disks object structure is documented below. A maximum of 59 disks can be mounted. Changing this creates a new
  instance.

* `tags` - (Optional, Map, ForceNew) Specifies the key/value pairs to associate with the instance. Changing this creates
  a new instance.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies a unique id in UUID format of enterprise project .
  Changing this creates a new instance.

* `agency_name` - (Optional, String, ForceNew) Specifies the IAM agency name which is created on IAM to provide
  temporary credentials for BMS to access cloud services. Changing this creates a new instance.

* `delete_eip_on_termination` - (Optional, Bool) Specifies whether the EIP is released when the instance is terminated.
  Defaults to *true*.

* `delete_disks_on_termination` - (Optional, Bool) Specifies whether to delete the data disks when the instance is
  terminated.
  Defaults to *false*.

The `nics` block supports:

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of subnet to attach to the instance. Changing this creates
  a new instance.

* `ip_address` - (Optional, String, ForceNew) Specifies a fixed IPv4 address to be used on this network. Changing this
  creates a new instance.

The `data_disks` block supports:

* `type` - (Required, String, ForceNew) Specifies the BMS data disk type, which must be one of available disk types.
  Changing this creates a new instance.

* `size` - (Required, Int, ForceNew) Specifies the data disk size, in GB. The value ranges form 10 to 32768. Changing
  this creates a new instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A resource ID in UUID format.
* `host_id` - The host ID of the instance.
* `status` - The status of the instance.
* `description` - The description of the instance.
* `image_name` - The image_name of the instance.
* `public_ip` - The EIP address that is associted to the instance.
* `nics` - An array of one or more networks to attach to the instance.
  The [nics_struct](#BMS_Response_nics_struct) structure is documented below.
* `disk_ids` - The ID of disks attached.

<a name="BMS_Response_nics_struct"></a>
The `nics_struct` block supports:

* `mac_address` - The MAC address of the nic.
* `port_id` - The port ID corresponding to the IP address.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minute.
* `update` - Default is 30 minute.
* `delete` - Default is 30 minute.
