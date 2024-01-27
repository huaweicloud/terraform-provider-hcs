---
subcategory: "Data Warehouse Service (DWS)"
---

# hcs_dws_cluster

Manages a DWS cluster resource within HuaweiCloudStack.  

## Example Usage

```hcl
variable "availability_zone" {}
variable "dws_cluster_name" {}
variable "dws_cluster_version" {}
variable "user_name" {}
variable "user_pwd" {}
variable "vpc_id" {}
variable "network_id" {}

resource "hcs_networking_secgroup" "secgroup" {
  name        = "sg_dws"
  description = "terraform security group"
}

resource "hcs_dws_cluster" "cluster" {
  name              = var.dws_cluster_name
  version           = var.dws_cluster_version
  node_type         = "dws.m3.xlarge"
  number_of_node    = 3
  availability_zone = var.availability_zone
  user_name         = var.user_name
  user_pwd          = var.user_pwd
  vpc_id            = var.vpc_id
  network_id        = var.network_id
  security_group_id = hcs_networking_secgroup.secgroup.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Cluster name, which must be unique and contains 4 to 64 characters, which
  consist of letters, digits, hyphens(-), or underscores(_) only and must start with a letter.
  Changing this creates a new cluster resource.

* `node_type` - (Required, String, ForceNew) The flavor of the cluster.  
  Changing this parameter will create a new resource.

* `number_of_node` - (Required, Int) Number of nodes in a cluster.  
  The value ranges from 3 to 32 in cluster mode. The value of stream warehouse(stand-alone mode) is 1.

* `user_name` - (Required, String, ForceNew) Administrator username for logging in to a data warehouse cluster.  
  The administrator username must: Consist of lowercase letters, digits, or underscores.
  Start with a lowercase letter or an underscore.
  Changing this parameter will create a new resource.

* `user_pwd` - (Required, String) Administrator password for logging in to a data warehouse cluster.
  A password contains 8 to 32 characters, which consist of letters, digits,
  and special characters(~!@#%^&*()-_=+|[{}];:,<.>/?).
  It cannot be the same as the username or the username written in reverse order.

* `vpc_id` - (Required, String, ForceNew) The VPC ID.
  Changing this parameter will create a new resource.

* `network_id` - (Required, String, ForceNew) The subnet ID.
  Changing this parameter will create a new resource.

* `security_group_id` - (Required, String, ForceNew) The security group ID.
  Changing this parameter will create a new resource.

* `availability_zone` - (Required, String, ForceNew) The availability zone in which to create the cluster instance.
  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project ID.
  Changing this parameter will create a new resource.

* `number_of_cn` - (Required, Int, ForceNew) The number of CN.  
  The value ranges from 2 to **number_of_node**, the maximum value is 20. Defaults to 3.
  Changing this parameter will create a new resource.

* `version` - (Required, String, ForceNew) The cluster version.
  Changing this parameter will create a new resource.

* `volume` - (Required, List, ForceNew) The information about the volume.

  Changing this parameter will create a new resource.

  The [Volume](#DwsCluster_Volume) structure is documented below.

* `port` - (Optional, Int, ForceNew) Service port of a cluster (8000 to 10000). The default value is 8000.  
  Changing this parameter will create a new resource.

* `tags` - (Optional, Map, ForceNew) The key/value pairs to associate with the cluster.
  Changing this parameter will create a new resource.

* `dss_pool_id` - (Optional, String, ForceNew) Dedicated storage pool ID.
  Changing this parameter will create a new resource.

* `kms_key_id` - (Optional, String, ForceNew) The KMS key ID.
  Changing this parameter will create a new resource.

* `public_ip` - (Optional, List, ForceNew) The information about public IP.  

  Changing this parameter will create a new resource.

  The [PublicIp](#DwsCluster_PublicIp) structure is documented below.

* `keep_last_manual_snapshot` - (Optional, Int) The number of latest manual snapshots that need to be
  retained when deleting the cluster.

<a name="DwsCluster_PublicIp"></a>
The `PublicIp` block supports:

* `public_bind_type` - (Optional, String) The bind type of public IP.  
  The valid value are **auto_assign**, **not_use**, and **bind_existing**. Defaults to **not_use**.

* `eip_id` - (Optional, String) The EIP ID.  

<a name="DwsCluster_Volume"></a>
The `Volume` block supports:

* `type` - (Optional, String) The volume type. Value options are as follows:
  + **SATA**: Common I/O. The SATA disk is used.
  + **SAS**: High I/O. The SAS disk is used.
  + **SSD**: Ultra-high I/O. The solid-state drive (SSD) is used.
  The valid value are **auto_assign**, **not_use**, and **bind_existing**. Defaults to **not_use**.

* `capacity` - (Optional, String) The capacity size, in GB.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `endpoints` - Private network connection information about the cluster.
  The [Endpoint](#DwsCluster_Endpoint) structure is documented below.

* `public_endpoints` - Public network connection information about the cluster.
  The [PublicEndpoint](#DwsCluster_PublicEndpoint) structure is documented below.

* `recent_event` - The recent event number.

* `status` - The cluster status.  
  The valid values are **CREATING**, **AVAILABLE**, **ACTIVE**, **FAILED**, **CREATE_FAILED**,
  **DELETING**, **DELETE_FAILED**, **DELETED**, and **FROZEN**.

* `sub_status` - Sub-status of clusters in the AVAILABLE state.  
  The value can be one of the following:
    + READONLY
    + REDISTRIBUTING
    + REDISTRIBUTION-FAILURE
    + UNBALANCED
    + UNBALANCED | READONLY
    + DEGRADED
    + DEGRADED | READONLY
    + DEGRADED | UNBALANCED
    + UNBALANCED | REDISTRIBUTING
    + UNBALANCED | REDISTRIBUTION-FAILURE
    + READONLY | REDISTRIBUTION-FAILURE
    + UNBALANCED | READONLY | REDISTRIBUTION-FAILURE
    + DEGRADED | REDISTRIBUTION-FAILURE
    + DEGRADED | UNBALANCED | REDISTRIBUTION-FAILURE
    + DEGRADED | UNBALANCED | READONLY | REDISTRIBUTION-FAILURE
    + DEGRADED | UNBALANCED | READONLY

* `task_status` - Cluster management task.  
  The value can be one of the following:
    + UNFREEZING
    + FREEZING
    + RESTORING
    + SNAPSHOTTING
    + GROWING
    + REBOOTING
    + SETTING_CONFIGURATION
    + CONFIGURING_EXT_DATASOURCE
    + DELETING_EXT_DATASOURCE
    + REBOOT_FAILURE
    + RESIZE_FAILURE

* `created` - The creation time of the cluster.  
  Format: ISO8601: **YYYY-MM-DDThh:mm:ssZ**.

* `updated` - The updated time of the cluster.  
  Format: ISO8601: **YYYY-MM-DDThh:mm:ssZ**.

* `private_ip` - List of private network IP addresses.  

* `maintain_window` - Cluster maintenance window.
  The [MaintainWindow](#DwsCluster_MaintainWindow) structure is documented below.

<a name="DwsCluster_Endpoint"></a>
The `Endpoint` block supports:

* `connect_info` - Private network connection information.

* `jdbc_url` - JDBC URL. Format: jdbc:postgresql://<connect_info>/<YOUR_DATABASE_NAME>  

<a name="DwsCluster_PublicEndpoint"></a>
The `PublicEndpoint` block supports:

* `public_connect_info` - Public network connection information.

* `jdbc_url` - JDBC URL. Format: jdbc:postgresql://<public_connect_info>/<YOUR_DATABASE_NAME>  

<a name="DwsCluster_MaintainWindow"></a>
The `MaintainWindow` block supports:

* `day` - Maintenance time in each week in the unit of day.  
  The valid values are **Mon**, **Tue**, **Wed**, **Thu**, **Fri**,
  **Sat**, and **Sun**.

* `start_time` - Maintenance start time in HH:mm format. The time zone is GMT+0.  

* `end_time` - Maintenance end time in HH:mm format. The time zone is GMT+0.  

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `update` - Default is 60 minutes.
* `delete` - Default is 60 minutes.

## Import

Cluster can be imported using the following format:

```
$ terraform import hcs_dws_cluster.test 47ad727e-9dcc-4833-bde0-bb298607c719
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `user_pwd`, `number_of_cn`, `kms_key_id`,
`volume`, `dss_pool_id`.
It is generally recommended running `terraform plan` after importing a cluster.
You can then decide if changes should be applied to the cluster, or the resource definition
should be updated to align with the cluster. Also you can ignore changes as below.

```
resource "hcs_dws_cluster" "test" {
    ...

  lifecycle {
    ignore_changes = [
      user_pwd, number_of_cn, kms_key_id, volume, dss_pool_id
    ]
  }
}
```
