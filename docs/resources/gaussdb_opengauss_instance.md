---
subcategory: "GaussDB"
---

# hcs_gaussdb_opengauss_instance

GaussDB OpenGauss instance management within HuaweiCouldStack.

-> **NOTE:** If the endpoint is manually configured, both **opengauss** and **opengaussv31** should be configured.

## Example Usage

### Create an instance for distributed HA mode

```hcl
variable "vpc_id" {}
variable "subnet_network_id" {}
variable "security_group_id" {}
variable "instance_name" {}
variable "instance_password" {}

data "hcs_availability_zones" "test" {}

resource "hcs_gaussdb_opengauss_instance" "test" {
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_network_id
  security_group_id = var.security_group_id

  flavor            = "gaussdb.opengauss.ee.dn.m6.2xlarge.8.in"
  name              = var.instance_name
  password          = var.instance_password
  sharding_num      = 1
  coordinator_num   = 2
  availability_zone = join(",", slice(data.hcs_availability_zones.test.names, 0, 3))

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "strong"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
```

### Create an instance for centralized HA mode

```hcl
variable "instance_name" {}
variable "instance_password" {}
variable "vpc_id" {}
variable "subnet_network_id" {}

data "hcs_availability_zones" "test" {}

resource "hcs_gaussdb_opengauss_instance" "instance_acc" {
  vpc_id            = var.vpc_id
  subnet_id         = var.subnet_network_id
  security_group_id = var.security_group_id
  name              = var.instance_name
  password          = var.instance_password
  flavor            = "gaussdb.opengauss.ee.m6.2xlarge.x868.ha"
  availability_zone = join(",", slice(data.hcs_availability_zones.myaz.names, 0, 3))

  replica_num = 3

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "strong"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the instance.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the instance name, which can be the same as an existing instance name.
  The value must be `4` to `64` characters in length and start with a letter. It is case-sensitive and can contain only
  letters, digits, hyphens (-), and underscores (_).

* `flavor` - (Required, String, ForceNew) Specifies the instance specifications. Please reference the API docs for valid
  options. Changing this parameter will create a new resource.

* `password` - (Required, String) Specifies the database password. The value must be `8` to `32` characters in length,
  including uppercase and lowercase letters, digits, and special characters, such as **~!@#%^*-_=+?**. You are advised
  to enter a strong password to improve security, preventing security risks such as brute force cracking.

* `availability_zone` - (Required, String, ForceNew) Specifies the availability zone information, can be three same or
  different az like **cn-north-4a,cn-north-4a,cn-north-4a**. Changing this parameter will create a new resource.

* `ha` - (Required, List, ForceNew) Specifies the HA information.
  The [ha](#opengauss_ha) structure is documented below.
  Changing this parameter will create a new resource.

* `volume` - (Required, List) Specifies the volume storage information.
  The [volume](#opengauss_volume) structure is documented below.

* `vpc_id` - (Required, String, ForceNew) Specifies the VPC ID to which the subnet belongs.
  Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the network ID of VPC subnet to which the instance belongs.
  Changing this parameter will create a new resource.

* `security_group_id` - (Optional, String, ForceNew) Specifies the security group ID to which the instance belongs.
  If the `port` parameter is specified, please ensure that the TCP ports in the inbound rule of security group
  includes the `100` ports starting with the database port.
  (For example, if the database port is `8,000`, the TCP port must include the range from `8,000` to `8,100`.)

  Changing this parameter will create a new resource.

* `port` - (Optional, String, ForceNew) Specifies the port information. Defaults to `8,000`.
  The valid values are as follows:
    + `2,378` to `2,380`
    + `4999` to `5,000`
    + `5,999` to `6,001`
    + `8,097` to `8,098`
    + `12,016` to `12,017`
    + `20,049` to `20,050`
    + `21,731` to `21,732`
    + `32,122` to `32,124`

  Changing this parameter will create a new resource.

* `configuration_id` - (Optional, String, ForceNew) Specifies the parameter template ID.
  Changing this parameter will create a new resource.

* `os_type` - (Optional, String, ForceNew) Specifies the OS type. The value is case sensitive. The GaussDB version
  matching the Hce operating system must be 8.102 or later. The valid values are **Hce**, **Euler** and **Kylin**.
  Default to **Hce**.

  Changing this parameter will create a new resource.

* `sharding_num` - (Optional, Int) Specifies the sharding number. The valid value is range form `1` to `9`.
  The default value is 3.

* `coordinator_num` - (Optional, Int) Specifies the coordinator number. Values: 1~9. The default value is 3.
  The value must not be greater than twice value of `sharding_num`.

* `replica_num` - (Optional, Int, ForceNew) The replica number. The valid values are **2** and **3**, defaults to **3**.
  Double replicas are only available for specific users and supports only instance versions are v1.3.0 or later.
  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID.
  Changing this parameter will create a new resource.

* `time_zone` - (Optional, String, ForceNew) Specifies the time zone. Defaults to **UTC+08:00**.
  Changing this parameter will create a new resource.

* `force_import` - (Optional, Bool) Specifies whether to import the instance with the given configuration instead of
  creation. If specified, try to import the instance instead of creation if the instance already existed.

* `datastore` - (Optional, List, ForceNew) Specifies the datastore information.
  The [datastore](#opengauss_datastore) structure is documented below.
  Changing this parameter will create a new resource.

* `backup_strategy` - (Optional, List) Specifies the advanced backup policy.
  The [backup_strategy](#opengauss_backup_strategy) structure is documented below.

<a name="opengauss_ha"></a>
The `ha` block supports:

* `mode` - (Required, String, ForceNew) Specifies the database mode.
  The valid values is **centralization_standard**.
  Changing this parameter will create a new resource.

* `replication_mode` - (Required, String, ForceNew) Specifies the database replication mode.
  Only **sync** is supported now. Changing this parameter will create a new resource.

* `consistency` - (Optional, String, ForceNew) Specifies the database consistency mode.
  The valid values are **strong** and **eventual**, not case sensitive.
  Changing this parameter will create a new resource.

<a name="opengauss_volume"></a>
The `volume` block supports:

* `type` - (Required, String, ForceNew) Specifies the volume type. Only **ULTRAHIGH** is supported now.
  Changing this parameter will create a new resource.

* `size` - (Required, Int) Specifies the volume size (in gigabytes). The valid value is range form `40` to `4,000`.
  - **ECS deployment scheme**: The value ranges from (Number of shards x 40 GB) to (Number of shards x 24 TB). The size
    must be an integer multiple of (Number of shards x 40). The upper limit of disk usage varies according to the CPU
    size.
  - **BMS Deployment Scheme**: This parameter is automatically calculated based on the selected flavor and cannot be
    specified. Even if this parameter is set, it does not take effect.

<a name="opengauss_datastore"></a>
The `datastore` block supports:

* `engine` - (Required, String, ForceNew) Specifies the database engine. Only **GaussDB(for openGauss)** is supported
  now. Changing this parameter will create a new resource.

* `version` - (Optional, String, ForceNew) Specifies the database version. Defaults to the latest version. Please
  reference to the API docs for valid options. Changing this parameter will create a new resource.

<a name="opengauss_backup_strategy"></a>
The `backup_strategy` block supports:

* `start_time` - (Required, String) Specifies the backup time window. Automated backups will be triggered during the
  backup time window. It must be a valid value in the **hh:mm-HH:MM** format. The current time is in the UTC format. The
  **HH** value must be `1` greater than the **hh** value. The values of mm and MM must be the same and must be set to
  **00**. Example value: **08:00-09:00**, **23:00-00:00**.

* `keep_days` - (Optional, Int) Specifies the number of days to retain the generated backup files. The value ranges from
  `0` to `732`. If this parameter is set to `0`, the automated backup policy is not set.
  If this parameter is not transferred, the automated backup policy is enabled by default.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the DB instance ID.

* `status` - Indicates the DB instance status.

* `type` - Indicates the database type.

* `private_ips` - Indicates the private IP address of the DB instance.

* `public_ips` - Indicates the public IP address of the DB instance.

* `endpoints` - Indicates the connection endpoints list of the DB instance. Example: [127.0.0.1:8000].

* `db_user_name` - Indicates the default username.

* `switch_strategy` - Indicates the switch strategy.

* `maintenance_window` - Indicates the maintenance window.

* `nodes` - Indicates the instance nodes information. The [nodes](#opengauss_nodes) object structure is
  documented below.

<a name="opengauss_nodes"></a>
The `nodes` block contains:

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `role` - Indicates the node role.
    + **master**.
    + **slave**.

* `status` - Indicates the node status.

* `availability_zone` - Indicates the availability zone of the node.

* `private_ip` - Indicates the internal IP address of the node. This parameter is valid only for CN nodes in the
  distributed edition. This parameter is valid for all nodes in the centralized edition.
  The parameter value exists after the ECS is created.

* `public_ip` - Indicates the bound external IP address of the node. This parameter is valid only for CN nodes in the
  distributed edition. This parameter is valid for all nodes in the centralized edition.
  The parameter value exists after the ECS is created and bound to an EIP.

* `data_ip` - Indicates the data ip of the node.

* `bms_hs_ip` - IP address of the high-speed NIC, which is a dedicated IP field of the BMS instance and is used for
  data synchronization.

* `management_ip` - Indicates the management ip of the node.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 120 minutes.
* `update` - Default is 90 minutes.
* `delete` - Default is 45 minutes.

## Import

OpenGaussDB instance can be imported using the `id`, e.g.

```
$ terraform import hcs_gaussdb_opengauss_instance.test 1f2c4f48adea4ae684c8edd8818fa349in14
```


Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include:
`password`, `availability_zone` and `ha.mode`.
It is generally recommended running `terraform plan` after importing a opengauss instance.
You can then decide if changes should be applied to the instance, or the resource
definition should be updated to align with the instance. Also, you can ignore changes as below.

```hcl
resource "hcs_gaussdb_opengauss_instance" "instance" {
  lifecycle {
    ignore_changes = [
      password, availability_zone,
    ]
  }
}
```
