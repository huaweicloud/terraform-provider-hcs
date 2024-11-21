---
subcategory: "Elastic Cloud Server (ECS)"
---

# hcs_ecs_compute_instance_clone

Clone a ECS VM instance resource within HCS.

## Example Usage

### Basic Instance

```hcl
resource "hcs_ecs_compute_instance_clone" "clone_basic" {
  instance_id = "c9a6edf4-ba26-48cb-857f-9aff9ffa25dd"
  power_on    = true
}
```

### Instance With Modify Network

```hcl
resource "hcs_ecs_compute_instance_clone" "clone_modify_network" {
  instance_id = "c9a6edf4-ba26-48cb-857f-9aff9ffa25dd"
  power_on    = true
  vpc_id = "ba69be72-406b-4d3a-b552-5f93fb47c1d9"
  network {
    fixed_ip_v4 = "ipv4 address"
    subnet_id = "b4db86e5-406e-422b-8744-54a15beb0aba"
    ipv6_enable = true
    fixed_ip_v6 = "ipv6 address"
    security_group_ids {
      security_group_id = "3e8cdc50-b6e3-4f7a-b7b7-88fe18a0e6a0"
    }
  }
}
```

### Instance With Modify Multiple Network

```hcl
resource "hcs_ecs_compute_instance_clone" "clone_modify_networks" {
  instance_id         = var.instance_id
  power_on          = var.power_on
  vpc_id = var.vpc_id

  network {
    subnet_id = "b4db86e5-406e-422b-8744-54a15beb0aba"
    fixed_ip_v4 = "192.167.0.10"
    ipv6_enable = false
    fixed_ip_v6 = "1030::C9B4:FF12:48AA:1A2B"
    security_group_ids {
      security_group_id = "3e8cdc50-b6e3-4f7a-b7b7-88fe18a0e6a0"
    }
  }

  network {
    subnet_id = "fa8b4105-7128-46ad-a44f-3f00549a6bb5"
    fixed_ip_v4 = "192.167.0.11"
    ipv6_enable = false
    fixed_ip_v6 = "1030::C9B4:FF12:48AA:1A2B"
    security_group_ids {
      security_group_id = "aa8d4c81-b53d-4db5-b66b-25b9d8c24896"
    }
    security_group_ids {
      security_group_id = "c8cdd497-1a66-4674-8687-d49df2deb40e"
    }
  }
}
```

### Instance With Modify Password

```hcl
resource "hcs_ecs_compute_instance_clone" "clone_modify_passwd" {
  instance_id = "c9a6edf4-ba26-48cb-857f-9aff9ffa25dd"
  retain_passwd = false
  power_on    = true
  admin_pass = "test_password"
}
```

### Instance With Modify Keypair

```hcl
resource "hcs_ecs_compute_instance_clone" "clone_modify_keypair" {
  instance_id = "c9a6edf4-ba26-48cb-857f-9aff9ffa25dd"
  retain_passwd = false
  power_on    = true
  key_pair = "keypair_name"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the Instance ID of the desired instance. Changing this clone a new instance.

* `name` - (Required, String, ForceNew) Specifies a unique name for the instance. The name consists of 1 to 64 characters,
  including letters, digits, underscores (_) and hyphens (-).

* `power_on` - (Required, bool, ForceNew) Specifies the power status of the instance.true: power-on; false: power-off.

* `retain_passwd` - (Optional, bool, ForceNew) Indicates whether to retain the old password. This parameter is valid only for Linux VMs.

* `admin_pass` - (Optional, String, ForceNew) Specifies the administrative password to assign to the instance, conflict with `key_pair`.

* `key_pair` - (Optional, String, ForceNew) Specifies the SSH keypair name used for logging in to the instance.

* `vpc_id` - (Optional, String, ForceNew) The VPC id of WAF dedicated instance. Changing this will create a new
  instance.

* `network` - (Optional, List, ForceNew) Specifies an array of one or more networks to attach to the instance. The network object structure is documented below. 

The `network` block supports:

* `subnet_id` - (Optional, String, ForceNew) The subnet id of WAF dedicated instance VPC. 

* `fixed_ip_v4` - (Optional, String, ForceNew) Specifies a fixed IPv4 address to be used on this network.

* `ipv6_enable` - (Optional, Bool, ForceNew) Specifies whether the IPv6 function is enabled for the nic.
  Defaults to false. 

* `fixed_ip_v6` - (Optional, String, ForceNew) Specifies a fixed IPv6 address to be used on this network.

* `security_group_ids` - (Optional, List, ForceNew) Specifies an array of one or more security group IDs to associate with the instance.

## Note

ECS resources can only be cloned, but cannot be updated or deleted locally.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minute.