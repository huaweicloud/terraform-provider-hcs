---
subcategory: "MapReduce Service (MRS)"
layout: "huaweicloudstack"
page_title: "HuaweiCloudStack: hcs_mrs_clusters"
description: |-
  Use this data source to get clusters of MapReduce.
---

# hcs_mrs_clusters

Use this data source to get clusters of MapReduce.

## Example Usage

```hcl
data "hcs_mrs_clusters" "test" {
  status = "running"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) The name of cluster.

* `status` - (Optional, String) The status of cluster.  
  The valid values are as follows:
  + **existing**: Query existing clusters, including all clusters except those in the deleted state
    and the yearly/monthly clusters in the Order processing or preparing state.
  + **history**: Quer historical clusters, including all the deleted clusters, clusters that fail to delete,
    clusters whose VMs fail to delete, and clusters whose database updates fail to delete.
  + **starting**: Query a list of clusters that are being started.
  + **running**: Query a list of running clusters.
  + **terminated**: Query a list of terminated clusters.
  + **failed**: Query a list of failed clusters.
  + **abnormal**: Query a list of abnormal clusters.
  + **terminating**: Query a list of clusters that are being terminated.
  + **frozen**: Query a list of frozen clusters.
  + **scaling-out**: Query a list of clusters that are being scaled out.
  + **scaling-in**: Query a list of clusters that are being scaled in.
  + **partial-error**: Query a list of some incorrect clusters.

* `enterprise_project_id` - (Optional, String) The enterprise project ID used to query clusters in a specified
  enterprise project.
  The default value is **0**, indicating the default enterprise project.

* `tags` - (Optional, String) You can search for a cluster by its tags.  
  If you specify multiple tags, the relationship between them is **AND**.
  The format of the tags parameter is **tags=k1\*v1,k2\*v2,k3\*v3**.
  When the values of some tags are null, the format is **tags=k1,k2,k3\*v3**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `clusters` - The list of clusters.
  The [clusters](#mrs_clusters_clusters_attr) structure is documented below.

<a name="mrs_clusters_clusters_attr"></a>
The `clusters` block supports:

* `id` - The ID of the cluster.

* `name` - The cluster name.

* `master_node_num` - The number of Master nodes deployed in a cluster.

* `core_node_num` - The number of Core nodes deployed in a cluster.

* `total_node_num` - The total number of nodes deployed in a cluster.

* `status` - The cluster status.  
  The following options are supported:  
  + **starting**: The cluster is being started.
  + **running**: The cluster is running.
  + **terminated**: The cluster has been terminated.
  + **failed**: The cluster fails.
  + **abnormal**: The cluster is abnormal.
  + **terminating**: The cluster is being terminated.
  + **frozen**: The cluster has been frozen.
  + **scaling-out**: The cluster is being scaled out.
  + **scaling-in**: The cluster is being scaled in.

* `billing_type` - The cluster billing mode.  
  The valid values are as follows:  
  + **11**: Yearly/Monthly.
  + **12**: Pay-per-use.

* `vpc_id` - The VPC ID.

* `subnet_id` - The subnet ID.

* `duration` - The cluster subscription duration.

* `fee` - The cluster creation fee, which is automatically calculated.

* `hadoop_version` - The Hadoop version.

* `master_node_size` - The instance specifications of a Master node.

* `core_node_size` - The instance specifications of a Core node.

* `component_list` - The component list.

  The [component_list](#mrs_clusters_component_list_attr) structure is documented below.

* `external_ip` - The external IP address.

* `external_alternate_ip` - The backup external IP address.

* `internal_ip` - The internal IP address.

* `deployment_id` - The cluster deployment ID.

* `description` - The cluster description.

* `order_id` - The cluster creation order ID.

* `master_node_product_id` - The product ID of a Master node.

* `master_node_spec_id` - The specification ID of a Master node.

* `core_node_product_id` - The product ID of a Core node.

* `core_node_spec_id` - The specification ID of a Core node.

* `availability_zone` - The AZ.

* `vnc` - The URI for remotely logging in to an ECS.

* `volume_size` - The disk storage space.

* `volume_type` - The disk type.

* `enterprise_project_id` - The enterprise project ID.

* `type` - The cluster type.  
  The valid values are as follows:  
  + **0**: analysis cluster.
  + **1**: streaming cluster.
  + **2**: hybrid cluster.
  + **3**: custom cluster.
  + **4**: Offline cluster.

* `security_group_id` - The security group ID.

* `slave_security_group_id` - The security group ID of a non-Master node.  
  Currently, one MRS cluster uses only one security group. Therefore, this field has been discarded.

* `stage_desc` - The cluster progress description.  
  The cluster **installation** progress includes:
  + **Verifying cluster parameters**: Cluster parameters are being verified.
  + **Applying for cluster resources**: Cluster resources are being applied for.
  + **Creating VMs**: The VMs are being created.
  + **Initializing VMs**: The VMs are being initialized.
  + **Installing MRS Manager**: MRS Manager is being installed.
  + **Deploying the cluster**: The cluster is being deployed.
  + **Cluster installation failed**: Failed to install the cluster.

  The cluster **scale-out** progress includes:
  + **Preparing for scale-out**: Cluster scale-out is being prepared.
  + **Creating VMs**: The VMs are being created.
  + **Initializing VMs**: The VMs are being initialized.
  + **Adding nodes to the cluster**: The nodes are being added to the cluster.
  + **Scale-out failed**: Failed to scale out the cluster.

  The cluster **scale-in** progress includes:
  + **Preparing for scale-in**: Cluster scale-in is being prepared.
  + **Decommissioning instance**: The instance is being decommissioned.
  + **Deleting VMs**: The VMs are being deleted.
  + **Deleting nodes from the cluster**: The nodes are being deleted from the cluster.
  + **Scale-in failed**: Failed to scale in the cluster.

  If the cluster installation, scale-out, or scale-in fails, stageDesc will display the failure cause.

* `safe_mode` - The running mode of an MRS cluster.  
  The valid values are as follows:
  + **0**: Normal cluster.
  + **1**: Security cluster.

* `version` - The cluster version.

* `node_public_cert_name` - The name of the key file.

* `master_node_ip` - The IP address of a Master node.

* `private_ip_first` - The preferred private IP address.

* `tags` - The tag information.

* `log_collection` - Whether to collect logs when cluster installation fails.  
  The valid values are as follows:
  + **0**: Do not collect logs.
  + **1**: Collect logs.

* `task_node_groups` - The list of Task nodes.
  The [node_groups](#mrs_clusters_node_groups_attr) structure is documented below.

* `node_groups` - The list of Master, Core and Task nodes.
  The [node_groups](#mrs_clusters_node_groups_attr) structure is documented below.

* `master_data_volume_type` - The data disk storage type of the Master node.  
  Currently, **SATA**, **SAS**, and **SSD** are supported.

* `master_data_volume_size` - The data disk storage space of the Master node.  
  To increase data storage capacity, you can add disks at the same time when creating a cluster.  
  The valid value is range from `100 GB` to `32,000 GB`.

* `master_data_volume_count` - The number of data disks of the Master node.  
  The value can be set to `1` only.

* `core_data_volume_type` - The data disk storage type of the Core node.  
  Currently, **SATA**, **SAS**, and **SSD** are supported.

* `core_data_volume_size` - The data disk storage space of the Core node.  
  To increase data storage capacity, you can add disks at the same time when creating a cluster.  
  The valid value is range from `100 GB` to `32,000 GB`.

* `core_data_volume_count` - The number of data disks of the Core node.

* `period_type` - The subscription type is yearly or monthly.  
  The following options are supported:  
  + **0**: monthly subscription.
  + **1**: yearly subscription.

* `scale` - The status of node changes.  
  If this parameter is left blank, no change operation is performed on a cluster node.  
  The options are as follows:
  + **Scaling-out**: The cluster is being scaled out.
  + **Scaling-in**: The cluster is being scaled in.
  + **scaling-error**: The cluster is in the running state and fails to be scaled in or out or the specifications
    fail to be scaled up for the last time.
  + **scaling-up**: The master node specifications are being scaled up.
  + **scaling_up_first**: The standby master node specifications are being scaled up.
  + **scaled_up_first**: The standby master node specifications have been scaled up.
  + **scaled-up-success**: The master node specifications have been scaled up.

* `eip_id` - The unique ID of the cluster EIP.

* `eip_address` - The IPv4 address of the cluster EIP.

* `eipv6_address` - The IPv6 address of the cluster EIP.  
  This parameter is not returned when an IPv4 address is used.

<a name="mrs_clusters_component_list_attr"></a>
The `component_list` block supports:

* `component_id` - The component ID.  
  For example, the component_id of Hadoop is MRS 3.0.2_001, MRS 2.1.0_001, MRS 1.9.2_001, MRS 1.8.10_001.

* `component_name` - The component name.

* `component_version` - The component version.

* `component_desc` - The component description.

<a name="mrs_clusters_node_groups_attr"></a>
The `node_groups` block supports:

* `group_name` - The node group name.

* `node_num` - The number of nodes in a node group.

* `node_size` - The instance specifications of a node group.

* `node_spec_id` - The instance specification ID of a node group.

* `node_product_id` - The instance product ID of a node group.

* `vm_product_id` - The VM product ID of a node group.

* `vm_spec_code` - The VM specification code of a node group.

* `root_volume_size` - The root disk storage space of a node group.

* `root_volume_type` - The root disk storage type of a node group.

* `root_volume_product_id` - The root disk product ID of a node group.

* `root_volume_resource_spec_code` - The root disk specification code of a node group.

* `root_volume_resource_type` - The system disk product type of a node group.

* `data_volume_type` - The data disk storage type of a node group.  
  The following options are supported:  
  + **SATA**: Common I/O.
  + **SAS**: High I/O.
  + **SSD**: Ultra-high I/O.

* `data_volume_count` - The number of data disks of a node group.

* `data_volume_size` - The data disk storage space of a node group.

* `data_volume_product_id` - The data disk product ID of a node group.

* `data_volume_resource_spec_code` - The data disk specification code of a node group.

* `data_volume_resource_type` - The data disk product type of a node group.
