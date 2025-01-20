---
subcategory: "SFS Turbo"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_sfs_turbo"
description: ""
---

# hcs_sfs_turbo

Provides a Shared File System (SFS) Turbo resource.

## Example Usage

### Create a hdd Shared File System (SFS) Turbo

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "secgroup_id" {}
variable "test_az" {}

resource "hcs_sfs_turbo" "test" {
  name              = "sfs-turbo-hdd"
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
  availability_zone = var.test_az

  share_proto = "NFS"
  share_type  = "sfsturbo.hdd"
  size        = 50
  bandwidth   = 150
  
  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### Create a ssd Shared File System (SFS) Turbo

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "secgroup_id" {}
variable "test_az" {}

resource "hcs_sfs_turbo" "test" {
  name              = "sfs-turbo-ssd"
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_id
  security_group_id = var.secgroup_id
  availability_zone = var.test_az

  share_proto = "NFS"
  share_type  = "sfsturbo.ssd"
  size        = 50
  bandwidth   = 350

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the SFS Turbo resource. If omitted, the
  provider-level region will be used. Changing this creates a new SFS Turbo resource.

* `name` - (Required, String) Specifies the name of an SFS Turbo file system. The value contains `4` to `64`
  characters and must start with a letter.

* `availability_zone` - (Required, String, ForceNew) Specifies the availability zone where the file system is located.

  Changing this will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID. Changing this will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the network ID of the subnet. 
  
  Changing this will create a new resource.

* `security_group_id` - (Required, String) Specifies the security group ID.

* `size` - (Required, Int) Specifies the capacity of a sharing file system, in GB. The valid range is from `50`
  to `1,048,576`.

* `bandwidth` - (Required, Int) Specifies the bandwidth of a sharing file system, in GB. 
  - When `share_type` is set to `sfsturbo.hdd`, the valid range of bandwidth is from `150` to `8192`.
  - When `share_type` is set to `sfsturbo.ssd`, the valid range of bandwidth is from `350` to `16384`.

  -> The file system capacity can only be expanded, not reduced.

* `share_proto` - (Optional, String, ForceNew) Specifies the protocol for sharing file systems. The valid value is
  **NFS**.

  Changing this will create a new resource.

* `share_type` - (Optional, String, ForceNew) Specifies the file system type. Changing this will create a new resource.
  Valid values are **sfsturbo.ssd** or **sfsturbo.hdd**.

  Changing this will create a new resource.

* `dedicated_flavor` - (Optional, String, ForceNew) Specifies the VM flavor used for creating a dedicated file system.

* `dedicated_storage_id` - (Optional, String, ForceNew) Specifies the ID of the dedicated distributed storage used
  when creating a dedicated file system.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project id of the file system. Changing this
  will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the SFS Turbo.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the SFS Turbo file system.

* `region` - The region of the SFS Turbo file system.

* `status` - The status of the SFS Turbo file system.

* `version` - The version ID of the SFS Turbo file system.

* `export_location` - The mount point of the SFS Turbo file system.

* `available_capacity` - The available capacity of the SFS Turbo file system in the unit of GB.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `update` - Default is 60 minutes.
* `delete` - Default is 10 minutes.

## Import

SFS Turbo can be imported using the `id`, e.g.

```bash
$ terraform import hcs_sfs_turbo 1e3d5306-24c9-4316-9185-70e9787d71ab
```
