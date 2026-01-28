---
subcategory: "MapReduce Service (MRS)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_mrs_cluster"
description: |-
  Use this data source to get cluster detail of MapReduce.
---

# hcs_mrs_cluster

Use this data source to get cluster detail of MapReduce.

## Example Usage

```hcl
variable "cluster_id" {}

data "hcs_mrs_cluster" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the cluster ID of MRS.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `cluster` - The detail of cluster.

  The [cluster](#mrs_cluster_cluster_attr) structure is documented below.

<a name="mrs_cluster_cluster_attr"></a>
The `cluster` block supports:

* `id` - The cluster ID.

* `cluster_name` - The cluster name.

* `master_node_num` - The number of master nodes deployed in a cluster.

* `core_node_num` - The number of core nodes deployed in a cluster.

* `total_node_num` - The total number of nodes deployed in a cluster.

* `cluster_state` - The status of cluster.
  + **starting**: The cluster is being started.
  + **running**: The cluster is running.
  + **terminated**: The cluster has been terminated.
  + **failed**: The cluster failed.
  + **abnormal**: The cluster is abnormal.
  + **terminating**: The cluster is being terminated.
  + **frozen**: The cluster has been frozen.
  + **scaling-out**: The cluster is being scaled out.
  + **scaling-in**: The cluster is being scaled in.

* `create_at` - The cluster creation time. such as 2025-07-22T15:32:36+08:00.

* `update_at` - The cluster update time. such as 2025-07-22T15:32:36+08:00.

* `billing_type` - The cluster billing mode.

* `data_center` - The cluster work region.

* `vpc_name` - The VPC name.

* `vpc_id` - The VPC ID.

* `duration` - The cluster subscription duration.

* `fee` - The cluster creation fee, which is automatically calculated.

* `hadoop_version` - The hadoop version.

* `master_node_size` - The instance specifications of a master node.

* `core_node_size` - The instance specifications of a core node.

* `component_list` - The component list of cluster.

  The [component_list](#cluster_component_list_struct) structure is documented below.

* `external_ip` - The external IP address.

* `external_alternate_ip` - The backup external IP address.

* `internal_ip` - The internal IP address.

* `deployment_id` - The cluster deployment ID.

* `remark` - The cluster remarks.

* `order_id` - The cluster creation order ID.

* `az_id` - The ID of availability zone.

* `master_node_product_id` - The product ID of a master node.

* `master_node_spec_id` - The specifications ID of a master node.

* `core_node_product_id` - The product ID of a core node.

* `core_node_spec_id` - The specifications ID of a core node.

* `az_name` - The name of availability zone.

* `az_code` - The name of availability zone(en).

* `availability_zone_id` - The availability zone ID.

* `instance_id` - The instance ID.

* `vnc` - The URI for remotely logging in to an ECS.

* `tenant_id` - The project ID.

* `volume_size` - The disk storage space.

* `volume_type` - The disk type.

* `subnet_id` - The subnet ID.

* `subnet_name` - The subnet name.

* `security_groups_id` - The security group ID.

* `slave_security_groups_id` - The security group ID of a non-master node. Currently, one MRS cluster uses only one
  security group. Therefore, this field has been **deprecated**. For compatibility purposes, the value of this field
  is the same as that of **security_groups_id**.

* `enterprise_project_id` - The enterprise project ID.

* `stage_desc` - The cluster operation progress description.
  + The cluster installation progress includes
    - **Verifying cluster parameters**: Cluster parameters are being verified.
    - **Applying for cluster resources**: Cluster resources are being requested.
    - **Creating VM**: The VMs are being created.
    - **Initializing VM**: The VMs are being initialized.
    - **Installing MRS Manager**: MRS Manager is being installed.
    - **Deploying cluster**: The cluster is being deployed.
    - **Cluster installation failed**: Failed to install the cluster.
  + The cluster scale-out progress includes:
    - **Preparing for cluster expansion**: Cluster scale-out is being prepared.
    - **Creating VM**: The VMs are being created.
    - **Initializing VM**: The VMs are being initialized.
    - **Adding node to the cluster**: The nodes are being added to the cluster.
    - **Cluster expansion failed**: Failed to scale out the cluster.
  + The cluster scale-in progress includes:
    - **Preparing for cluster shrink**: Cluster scale-in is being prepared.
    - **Decommissioning instance**: The instance is being decommissioned.
    - **Deleting VM**: The VMs are being deleted.
    - **Deleting node from the cluster**: The nodes are being deleted from the cluster.
    - **Cluster shrink failed**: Failed to scale in the cluster.

* `is_mrs_manager_finish` - Whether MRS Manager installation is complete during cluster creation.
  + true: MRS Manager installation is complete.
  + false: MRS Manager installation is not complete.

* `safe_mode` - The running mode of an MRS cluster.
  + **0**: Normal cluster.
  + **1**: Security cluster.

* `cluster_version` - The cluster version.

* `node_public_cert_name` - The name of the public key file.

* `master_node_ip` - The IP address of a master node.

* `private_ip_first` - The preferred private IP address.

* `error_info` - The error information.

* `tags` - The tag information.

* `charging_start_time` - The start time of billing.

* `cluster_type` - The cluster type.

* `log_collection` - Whether to collect logs when cluster installation fails.
  + **0**: Do not collect.
  + **1**: Collect.

* `task_node_groups` - The list of task nodes of cluster.

  The [task_node_groups](#cluster_node_groups_struct) structure is documented below.

* `node_groups` - The list of master, core, and task nodes.

  The [node_groups](#cluster_node_groups_struct) structure is documented below.

* `master_data_volume_type` - The data disk storage type of the master node.
  Currently, **SATA**, **SAS** and **SSD** are supported.

* `master_data_volume_size` - The data disk storage space of the master node.

* `master_data_volume_count` - The number of data disks of the master node.

* `core_data_volume_type` - The data disk storage type of the core node.
  Currently, **SATA**, **SAS** and **SSD** are supported.

* `core_data_volume_size` - The data disk storage space of the core node.

* `core_data_volume_count` - The number of data disks of the core node.

* `period_type` - Whether the subscription type is yearly or monthly.
  + `0`: monthly subscription.
  + `1`: yearly subscription.

* `scale` - The node change status. If this parameter is left blank, the cluster nodes are not changed.
  + **scaling-out**: The cluster is being scaled out.
  + **scaling-in**: The cluster is being scaled in.
  + **scaling-error**: The cluster is in the running state and fails in the last scale-out, scale-in, or specifications upgrade.
  + **scaling-up**: The master node specifications are being scaled up.
  + **scaling_up_first**: The standby master node specifications are being scaled up.
  + **scaled_up_first**: The standby master node specifications have been scaled up successfully.
  + **scaled-up-success**: The master node specifications have been scaled up successfully.

* `oms_business_ip` - The cluster OMS Master Node Business IP.
  It returned only for clusters managed by MRS **3.1.3-LTS** and **later** versions.

* `oms_alternate_business_ip` - The cluster OMS Standby Node Business IP.
  It returned only for clusters managed by MRS **3.1.3-LTS** and **later** versions.

* `oms_business_ip_port` - The Port bound to the service IP addresses of the active and standby OMS in the cluster.
  It returned only for clusters managed by MRS **3.1.3-LTS** and **later** versions.

<a name="cluster_component_list_struct"></a>
The `component_list` block supports:

* `component_id` - The component ID.

* `component_name` - The component name.

* `component_version` - The component version.

* `component_desc` - The component description.

<a name="cluster_node_groups_struct"></a>
The `node_groups` and `task_node_groups` block supports:

* `group_name` - The node group name.

* `node_num` - The number of nodes. The value ranges from 0 to 500. The minimum number of master and core nodes
  is 1 and the total number of core and task nodes cannot exceed 500.

* `node_size` - The instance specifications of a node.

* `node_spec_id` - The instance specifications ID of a node.

* `node_product_id` - The instance product ID of a node.

* `vm_product_id` - The VM product ID of a node.

* `vm_spec_code` - The VM specifications of a node.

* `root_volume_size` - The system disk size of a node.

* `root_volume_product_id` - The system disk product ID of a node.

* `root_volume_type` - The system disk type of a node.

* `root_volume_resource_spec_code` - The system disk product specifications of a node.

* `root_volume_resource_type` - The system disk product type of a node.。

* `data_volume_type` - The data disk storage type of a node.
  + **SATA**: Common I/O.
  + **SAS**: High I/O.
  + **SSD**: Ultra-high I/O.

* `data_volume_count` - The number of data disks of a node.

* `data_volume_size` - The data disk storage space of a node.

* `data_volume_product_id` - The data disk product ID of a node.

* `data_volume_resource_spec_code` - The data disk product specifications of a node.

* `data_volume_resource_type` - The data disk product type of a node.
