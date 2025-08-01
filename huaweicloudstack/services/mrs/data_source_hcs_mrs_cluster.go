package mrs

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

// MRS v1.1/{project_id}/cluster_infos/{cluster_id}
func DataSourceMrsCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMrsMapreduceClusterRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the resource.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The cluster ID of MRS`,
			},
			"cluster": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The detail of cluster.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cluster ID.`,
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cluster name.`,
						},
						"master_node_num": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The number of master nodes deployed in a cluster.`,
						},
						"core_node_num": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The number of core nodes deployed in a cluster.`,
						},
						"total_node_num": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The total number of nodes deployed in a cluster.`,
						},
						"cluster_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of cluster.`,
						},
						"create_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cluster creation time.`,
						},
						"update_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cluster update time.`,
						},
						"billing_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cluster billing mode.`,
						},
						"data_center": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cluster work region.`,
						},
						"vpc_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The VPC name.`,
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The VPC ID.`,
						},
						"duration": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cluster subscription duration.`,
						},
						"fee": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cluster creation fee, which is automatically calculated.`,
						},
						"hadoop_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The hadoop version.`,
						},
						"master_node_size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The instance specifications of a master node.`,
						},
						"core_node_size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The instance specifications of a core node.`,
						},
						"component_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The component list of cluster.`,
							Elem:        mrsClusterComponentListSchema(),
						},
						"external_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The external IP address.`,
						},
						"external_alternate_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The backup external IP address.`,
						},
						"internal_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The internal IP address.`,
						},
						"deployment_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cluster deployment ID.`,
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cluster remarks.`,
						},
						"order_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cluster creation order ID.`,
						},
						"az_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of availability zone.`,
						},
						"master_node_product_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The product ID of a master node.`,
						},
						"master_node_spec_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The specifications ID of a master node.`,
						},
						"core_node_product_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The product ID of a core node.`,
						},
						"core_node_spec_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The specifications ID of a core node.`,
						},
						"az_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of availability zone.`,
						},
						"az_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of availability zone(en).`,
						},
						"availability_zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The availability zone ID.`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The instance ID.`,
						},
						"vnc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The URI for remotely logging in to an ECS.`,
						},
						"tenant_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The project ID.`,
						},
						"volume_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The disk storage space.`,
						},
						"volume_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The disk type.`,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The subnet ID.`,
						},
						"subnet_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The subnet name.`,
						},
						"security_groups_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The security group ID.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The enterprise project ID.`,
						},
						"slave_security_groups_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The security group ID of a non-master node.`,
						},
						"stage_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cluster operation progress description.`,
						},
						"is_mrs_manager_finish": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether MRS Manager installation is complete during cluster creation.`,
						},
						"safe_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The running mode of an MRS cluster.`,
						},
						"cluster_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cluster version.`,
						},
						"node_public_cert_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the public key file.`,
						},
						"master_node_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The IP address of a master node.`,
						},
						"private_ip_first": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The preferred private IP address.`,
						},
						"error_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The error information.`,
						},
						"tags": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The tag information.`,
						},
						"charging_start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The start time of billing.`,
						},
						"cluster_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The cluster type.`,
						},
						"log_collection": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Whether to collect logs when cluster installation fails.`,
						},
						"task_node_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of task nodes of cluster.`,
							Elem:        mrsClusterNodeGroupSchema(),
						},
						"node_groups": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of master, core, and task nodes.`,
							Elem:        mrsClusterNodeGroupSchema(),
						},
						"master_data_volume_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The data disk storage type of the master node.`,
						},
						"master_data_volume_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The data disk storage space of the master node.`,
						},
						"master_data_volume_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of data disks of the master node.`,
						},
						"core_data_volume_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The data disk storage type of the core node.`,
						},
						"core_data_volume_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The data disk storage space of the core node.`,
						},
						"core_data_volume_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of data disks of the core node.`,
						},
						"period_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Whether the subscription type is yearly or monthly.`,
						},
						"scale": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The node change status. If this parameter is left blank, the cluster nodes are not changed.`,
						},

						// unique field of HCS
						"oms_business_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: ` The cluster OMS Master Node Business IP.`,
						},
						"oms_alternate_business_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cluster OMS Standby Node Business IP.`,
						},
						"oms_business_ip_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The Port bound to the service IP addresses of the active and standby OMS in the cluster.`,
						},
					},
				},
			},
		},
	}
}

// taskNodeGroup and nodeGroups
func mrsClusterNodeGroupSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The node group name.`,
			},
			"node_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of nodes. The value ranges from 0 to 500. The minimum number of master and core nodes is 1 and the total number of core and task nodes cannot exceed 500.`,
			},
			"node_size": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance specifications of a node.`,
			},
			"node_spec_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance specifications ID of a node.`,
			},
			"node_product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance product ID of a node.`,
			},
			"vm_product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The VM product ID of a node.`,
			},
			"vm_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The VM specifications of a node.`,
			},
			"root_volume_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The system disk size of a node. This parameter is not configurable and its default value is **40 GB**.`,
			},
			"root_volume_product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The system disk product ID of a node.`,
			},
			"root_volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The system disk type of a node.`,
			},
			"root_volume_resource_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The system disk product specifications of a node.`,
			},
			"root_volume_resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The system disk product type of a node.。`,
			},
			"data_volume_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Data disk storage type of a node. Currently, SATA, SAS, and SSD are supported.`,
			},
			"data_volume_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of data disks of a node.`,
			},
			"data_volume_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The data disk storage space of a node.`,
			},
			"data_volume_product_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data disk product ID of a node.`,
			},
			"data_volume_resource_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data disk product specifications of a node.`,
			},
			"data_volume_resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data disk product type of a node.`,
			},
		},
	}

	return &sc
}

func mrsClusterComponentListSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"component_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The component ID.`,
			},
			"component_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The component name.`,
			},
			"component_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The component version.`,
			},
			"component_desc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The component description.`,
			},
		},
	}

	return &sc
}

func getClusterDetail(client *golangsdk.ServiceClient, clusterId string) (interface{}, error) {
	httpUrl := "v1.1/{project_id}/cluster_infos/{cluster_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", clusterId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("cluster", respBody, nil), nil
}

func dataSourceMrsMapreduceClusterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	client, err := cfg.MrsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating MRS V1 client: %s", err)
	}

	clusterId := d.Get("cluster_id").(string)
	resp, err := getClusterDetail(client, clusterId)
	if err != nil {
		return diag.Errorf("error getting MRS cluster detail: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	cluster, err := flattenCluster(resp)
	if err != nil {
		return diag.Errorf("unable to flatten cluster: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("cluster", cluster),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCluster(resp interface{}) ([]map[string]interface{}, error) {
	if resp == nil {
		return nil, nil
	}

	createStr := utils.PathSearch("createAt", resp, false).(string)
	updateStr := utils.PathSearch("updateAt", resp, false).(string)
	createAt, err := strconv.ParseInt(createStr, 10, 64)
	if err != nil {
		return nil, err
	}
	updateAt, err := strconv.ParseInt(updateStr, 10, 64)
	if err != nil {
		return nil, err
	}

	return []map[string]interface{}{
		{
			"id":                        utils.PathSearch("clusterId", resp, nil),
			"cluster_name":              utils.PathSearch("clusterName", resp, nil),
			"master_node_num":           utils.PathSearch("masterNodeNum", resp, nil),
			"core_node_num":             utils.PathSearch("coreNodeNum", resp, nil),
			"total_node_num":            utils.PathSearch("totalNodeNum", resp, nil),
			"cluster_state":             utils.PathSearch("clusterState", resp, nil),
			"billing_type":              utils.PathSearch("billingType", resp, nil),
			"data_center":               utils.PathSearch("dataCenter", resp, nil),
			"vpc_name":                  utils.PathSearch("vpc", resp, nil),
			"vpc_id":                    utils.PathSearch("vpcId", resp, nil),
			"duration":                  utils.PathSearch("duration", resp, nil),
			"fee":                       utils.PathSearch("fee", resp, nil),
			"hadoop_version":            utils.PathSearch("hadoopVersion", resp, nil),
			"master_node_size":          utils.PathSearch("masterNodeSize", resp, nil),
			"core_node_size":            utils.PathSearch("coreNodeSize", resp, nil),
			"component_list":            flattenClustersComponentList(resp),
			"external_ip":               utils.PathSearch("externalIp", resp, nil),
			"external_alternate_ip":     utils.PathSearch("externalAlternateIp", resp, nil),
			"internal_ip":               utils.PathSearch("internalIp", resp, nil),
			"deployment_id":             utils.PathSearch("deploymentId", resp, nil),
			"remark":                    utils.PathSearch("remark", resp, nil),
			"order_id":                  utils.PathSearch("orderId", resp, nil),
			"az_id":                     utils.PathSearch("azId", resp, nil),
			"master_node_product_id":    utils.PathSearch("masterNodeProductId", resp, nil),
			"master_node_spec_id":       utils.PathSearch("masterNodeSpecId", resp, nil),
			"core_node_product_id":      utils.PathSearch("coreNodeProductId", resp, nil),
			"core_node_spec_id":         utils.PathSearch("coreNodeSpecId", resp, nil),
			"az_name":                   utils.PathSearch("azName", resp, nil),
			"az_code":                   utils.PathSearch("azCode", resp, nil),
			"availability_zone_id":      utils.PathSearch("availabilityZoneId", resp, nil),
			"instance_id":               utils.PathSearch("instanceId", resp, nil),
			"vnc":                       utils.PathSearch("vnc", resp, nil),
			"tenant_id":                 utils.PathSearch("tenantId", resp, nil),
			"volume_size":               utils.PathSearch("volumeSize", resp, nil),
			"volume_type":               utils.PathSearch("volumeType", resp, nil),
			"subnet_id":                 utils.PathSearch("subnetId", resp, nil),
			"subnet_name":               utils.PathSearch("subnetName", resp, nil),
			"security_groups_id":        utils.PathSearch("securityGroupsId", resp, nil),
			"enterprise_project_id":     utils.PathSearch("enterpriseProjectId", resp, nil),
			"slave_security_groups_id":  utils.PathSearch("slaveSecurityGroupsId", resp, nil),
			"stage_desc":                utils.PathSearch("stageDesc", resp, nil),
			"is_mrs_manager_finish":     utils.PathSearch("isMrsManagerFinish", resp, nil),
			"safe_mode":                 utils.PathSearch("safeMode", resp, nil),
			"cluster_version":           utils.PathSearch("clusterVersion", resp, nil),
			"node_public_cert_name":     utils.PathSearch("nodePublicCertName", resp, nil),
			"master_node_ip":            utils.PathSearch("masterNodeIp", resp, nil),
			"private_ip_first":          utils.PathSearch("privateIpFirst", resp, nil),
			"error_info":                utils.PathSearch("errorInfo", resp, nil),
			"tags":                      utils.PathSearch("tags", resp, nil),
			"charging_start_time":       utils.PathSearch("chargingStartTime", resp, nil),
			"cluster_type":              utils.PathSearch("clusterType", resp, nil),
			"log_collection":            utils.PathSearch("logCollection", resp, nil),
			"task_node_groups":          flattenClustersNodeGroups(resp, "taskNodeGroups"),
			"node_groups":               flattenClustersNodeGroups(resp, "nodeGroups"),
			"master_data_volume_type":   utils.PathSearch("masterDataVolumeType", resp, nil),
			"master_data_volume_size":   utils.PathSearch("masterDataVolumeSize", resp, nil),
			"master_data_volume_count":  utils.PathSearch("masterDataVolumeCount", resp, nil),
			"core_data_volume_type":     utils.PathSearch("coreDataVolumeType", resp, nil),
			"core_data_volume_size":     utils.PathSearch("coreDataVolumeSize", resp, nil),
			"core_data_volume_count":    utils.PathSearch("coreDataVolumeCount", resp, nil),
			"period_type":               utils.PathSearch("periodType", resp, nil),
			"scale":                     utils.PathSearch("scale", resp, nil),
			"oms_business_ip":           utils.PathSearch("omsBusinessIp", resp, nil),
			"oms_alternate_business_ip": utils.PathSearch("omsAlternateBusinessIp", resp, nil),
			"oms_business_ip_port":      utils.PathSearch("omsBusinessIpPort", resp, nil),
			"create_at":                 utils.FormatTimeStampRFC3339(createAt, false),
			"update_at":                 utils.FormatTimeStampRFC3339(updateAt, false),
		},
	}, nil
}

func flattenClustersComponentList(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("componentList", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"component_id":      utils.PathSearch("componentId", v, nil),
			"component_name":    utils.PathSearch("componentName", v, nil),
			"component_version": utils.PathSearch("componentVersion", v, nil),
			"component_desc":    utils.PathSearch("componentDesc", v, nil),
		})
	}

	return rst
}

func flattenClustersNodeGroups(resp interface{}, key string) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch(key, resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"group_name":                     utils.PathSearch("GroupName", v, nil),
			"node_num":                       utils.PathSearch("NodeNum", v, nil),
			"node_size":                      utils.PathSearch("NodeSize", v, nil),
			"node_spec_id":                   utils.PathSearch("NodeSpecId", v, nil),
			"node_product_id":                utils.PathSearch("NodeProductId", v, nil),
			"vm_product_id":                  utils.PathSearch("VmProductId", v, nil),
			"vm_spec_code":                   utils.PathSearch("VmSpecCode", v, nil),
			"root_volume_size":               utils.PathSearch("RootVolumeSize", v, nil),
			"root_volume_product_id":         utils.PathSearch("RootVolumeProductId", v, nil),
			"root_volume_type":               utils.PathSearch("RootVolumeType", v, nil),
			"root_volume_resource_spec_code": utils.PathSearch("RootVolumeResourceSpecCode", v, nil),
			"root_volume_resource_type":      utils.PathSearch("RootVolumeResourceType", v, nil),
			"data_volume_type":               utils.PathSearch("DataVolumeType", v, nil),
			"data_volume_count":              utils.PathSearch("DataVolumeCount", v, nil),
			"data_volume_size":               utils.PathSearch("DataVolumeSize", v, nil),
			"data_volume_product_id":         utils.PathSearch("DataVolumeProductId", v, nil),
			"data_volume_resource_spec_code": utils.PathSearch("DataVolumeResourceSpecCode", v, nil),
			"data_volume_resource_type":      utils.PathSearch("DataVolumeResourceType", v, nil),
		})
	}

	return rst
}
