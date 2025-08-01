package romaconnect

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/romaconnect/v2/instances"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/services/vpc"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

const (
	PendingStatus      = "PENDING"
	CreatingStatus     = "CREATING"
	RunningStatus      = "RUNNING"
	CreateFailedStatus = "CREATE_FAILED"
	ErrorStatus        = "ERROR"
	DeletedStatus      = "DELETED"
)

// ROMA Connect POST /v2/{project_id}/instances
// ROMA Connect GET /v2/{project_id}/{instance_id}
// ROMA Connect GET /v2/{project_id}/instances/{instance_id}/process
// ROMA Connect DELETE /v2/{project_id}/roma/instances/{instance_id}
// ROMA Connect GET /v2/{project_id}/instances?{query}
// VPC GET /v1/{project_id}/vpcs/{vpc_id}
// VPC GET /v1/{project_id}/subnets/{subnet_id}
func ResourceRomaConnectInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceCreate,
		ReadContext:   resourceInstanceRead,
		DeleteContext: resourceInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"available_zones": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ipv6_enable": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"enable_all": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"eip_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"entrance_bandwidth_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"publicip_enable"},
			},
			"mqs": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine_version": {
							Type:          schema.TypeString,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"rocketmq_enable"},
						},
						"rocketmq_enable": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"retention_policy": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"enable_publicip": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"ssl_enable": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"trace_enable": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"vpc_client_plain": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"connector_enable": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"enable_acl": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"cpu_architecture": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"maintain_begin": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				RequiredWith: []string{"maintain_end"},
			},
			"maintain_end": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"site_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_zone_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_arch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publicip_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publicip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publicip_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"connect_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charge_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidths": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mqs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enable": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"retention_policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ssl_enable": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"trace_enable": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"vpc_client_plain": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"partition_num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"specification": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_connect_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"public_connect_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"private_restful_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"public_restful_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"ipv6_connect_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rocketmq_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"external_elb_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"external_elb_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_elb_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_eip_bound": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_eip_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_eip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := config.GetHcsConfig(meta)
	region := conf.GetRegion(d)
	romaConnectClientV1, err := conf.RomaConnectV1Client(region)
	if err != nil {
		return diag.Errorf("[Create]Error creating ROMA Connect v1 client : %s", err)
	}
	romaConnectClientV2, err := conf.RomaConnectV2Client(region)
	if err != nil {
		return diag.Errorf("[Create]Error creating ROMA Connect v2 client: %s", err)
	}

	// check network
	vpcId := d.Get("vpc_id").(string)
	_, err = vpc.GetVpcById(conf, region, vpcId)
	if err != nil {
		return diag.Errorf("unable to find the vpc (%s) on the server: %s", vpcId, err)
	}

	subnetId := d.Get("subnet_id").(string)
	_, err = vpc.GetVpcSubnetById(conf, region, subnetId)
	if err != nil {
		return diag.Errorf("unable to find the subnet (%s) on the server: %s", subnetId, err)
	}

	// the available_zones of API is []string
	var az []string
	if availableZones, ok := d.GetOk("available_zones"); ok {
		az = utils.ExpandToStringList(availableZones.([]interface{}))
	}

	createOpts := instances.CreateOpts{
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		ProductId:             d.Get("product_id").(string),
		AvailableZones:        az,
		VpcId:                 vpcId,
		SubnetId:              subnetId,
		SecurityGroupId:       d.Get("security_group_id").(string),
		EnterpriseProjectId:   d.Get("enterprise_project_id").(string),
		Ipv6Enable:            utils.Bool(d.Get("ipv6_enable").(bool)),
		EnableAll:             utils.Bool(d.Get("enable_all").(bool)),
		EipId:                 d.Get("eip_id").(string),
		EntranceBandwidthSize: d.Get("entrance_bandwidth_size").(int),
		MaintainBegin:         d.Get("maintain_begin").(string),
		MaintainEnd:           d.Get("maintain_end").(string),
		CpuArchitecture:       d.Get("cpu_architecture").(string),
		Mqs:                   buildResourceRomaConnectInstanceMqs(d),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	n, err := instances.Create(romaConnectClientV1, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating ROMA Connect Instance: %s", err)
	}

	log.Printf("[INFO] Waiting for ROMA Connect Instance(%s) to become available", n.ID)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{CreatingStatus},
		Target:       []string{RunningStatus},
		Refresh:      waitForInstanceRunning(romaConnectClientV2, n.ID, []string{RunningStatus}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        15 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return diag.Errorf("Error waiting for ROMA Connect Instance(%s) to become available: %s", n.ID, stateErr)
	}

	d.SetId(n.ID)

	return resourceInstanceRead(ctx, d, meta)
}

func getRomaInstanceStatus(client *golangsdk.ServiceClient, d *schema.ResourceData) diag.Diagnostics {
	n, err := instances.GetProcess(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "[GET Process] ROMA Connect Instance")
	}

	// The values are 'CREATING', 'RUNNING' or 'CREATE_FAILED'
	if n.Instance.Status == RunningStatus {
		return nil
	}
	return diag.Errorf("The ROMA instance is not running: %v", n.Instance)
}

func resourceInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := config.GetHcsConfig(meta)
	region := conf.GetRegion(d)
	romaConnectClient, err := conf.RomaConnectV2Client(region)
	if err != nil {
		return diag.Errorf("[Read]Error creating ROMA Connect client: %s", err)
	}

	if err := getRomaInstanceStatus(romaConnectClient, d); err != nil {
		return err
	}

	n, err := instances.Get(romaConnectClient, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("[Read]Error get ROMA instance detail: %s", err)
	}

	log.Printf("[DEBUG] retrieving ROMA Connect Instance: %#v", n)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", n.Name),
		d.Set("site_id", n.SiteId),
		d.Set("description", n.Description),
		d.Set("flavor_id", n.FlavorId),
		d.Set("flavor_type", n.FlavorType),
		d.Set("project_id", n.ProjectId),
		d.Set("available_zone_ids", n.AvailableZoneIds),
		d.Set("vpc_id", n.VpcId),
		d.Set("subnet_id", n.SubnetId),
		d.Set("security_group_id", n.SecurityGroupId),
		d.Set("cpu_arch", n.CpuArch),
		d.Set("status", n.Status),
		d.Set("publicip_id", n.PublicIpId),
		d.Set("publicip_address", n.PublicIpAddress),
		d.Set("publicip_enable", n.PublicIpEnable),
		d.Set("connect_address", n.ConnectAddress),
		d.Set("charge_type", n.ChargeType),
		d.Set("bandwidths", n.Bandwidths),
		d.Set("ipv6_enable", n.Ipv6Enable),
		d.Set("maintain_begin", n.MaintainBegin),
		d.Set("maintain_end", n.MaintainEnd),
		d.Set("enterprise_project_id", n.EnterpriseProjectId),
		d.Set("resources", flattenResourceRomaResource(n.Resources)),
		d.Set("ipv6_connect_address", n.Ipv6ConnectAddress),
		d.Set("rocketmq_enable", n.RocketmqEnable),
		d.Set("external_elb_enable", n.ExternalElbEnable),
		d.Set("external_elb_id", n.ExternalElbId),
		d.Set("external_elb_address", n.ExternalElbAddress),
		d.Set("external_eip_bound", n.ExternalEipBound),
		d.Set("external_eip_id", n.ExternalEipId),
		d.Set("external_eip_address", n.ExternalEipAddress),
		d.Set("create_time", n.CreateTime),
		d.Set("update_time", n.UpdateTime),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := config.GetHcsConfig(meta)
	region := conf.GetRegion(d)
	romaConnectClientV1, err := conf.RomaConnectV1Client(region)
	if err != nil {
		return diag.Errorf("[Delete]Error creating ROMA Connect v1 client: %s", err)
	}
	romaConnectClientV2, err := conf.RomaConnectV2Client(region)
	if err != nil {
		return diag.Errorf("[Delete]Error creating ROMA Connect v2 client: %s", err)
	}

	err = instances.Delete(romaConnectClientV1, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("Error deleting ROMA Connect instance %s: %s", d.Id(), err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{PendingStatus},
		Target:       []string{DeletedStatus},
		Refresh:      waitForInstanceDeleted(romaConnectClientV2, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        15 * time.Second,
		PollInterval: 30 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error deleting ROMA Connect instance %s: %s", d.Id(), err)
	}
	return nil
}

func buildResourceRomaConnectInstanceMqs(d *schema.ResourceData) instances.MqsOpts {
	var mqsOpts instances.MqsOpts
	mqsRaw := d.Get("mqs").([]interface{})
	if len(mqsRaw) != 1 {
		return mqsOpts
	}
	mqs := mqsRaw[0].(map[string]interface{})

	// kafka and rocketmq common
	mqsOpts.SslEnable = mqs["ssl_enable"].(bool)
	mqsOpts.EnablePublicIp = mqs["enable_publicip"].(bool)

	// create a rocketmq instance
	if mqs["rocketmq_enable"].(bool) {
		mqsOpts.RocketMqEnable = mqs["rocketmq_enable"].(bool)
		mqsOpts.EnableAcl = mqs["enable_acl"].(bool)
		return mqsOpts
	}

	// create a kafka instance
	mqsOpts.EngineVersion = mqs["engine_version"].(string)
	mqsOpts.RetentionPolicy = mqs["retention_policy"].(string)
	mqsOpts.TraceEnable = mqs["trace_enable"].(bool)
	mqsOpts.VpcClientPlain = mqs["vpc_client_plain"].(bool)
	mqsOpts.ConnectorEnable = mqs["connector_enable"].(bool)

	return mqsOpts
}

func flattenResourceRomaResource(resource instances.Resources) []map[string]interface{} {
	mqs := resource.Mqs
	resourceData := []map[string]interface{}{
		{
			"mqs": []map[string]interface{}{
				{
					"id":                      mqs.ID,
					"enable":                  mqs.Enable,
					"retention_policy":        mqs.RetentionPolicy,
					"ssl_enable":              mqs.SslEnable,
					"trace_enable":            mqs.TraceEnable,
					"vpc_client_plain":        mqs.VpcClientPlain,
					"partition_num":           mqs.PartitionNum,
					"specification":           mqs.Specification,
					"private_connect_address": mqs.PrivateConnectAddress,
					"public_connect_address":  mqs.PublicConnectAddress,
					"private_restful_address": mqs.PrivateRestfulAddress,
					"public_restful_address":  mqs.PublicRestfulAddress,
				},
			},
		},
	}

	return resourceData
}

func waitForInstanceRunning(client *golangsdk.ServiceClient, instanceId string,
	targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := instances.GetProcess(client, instanceId).Extract()
		if err != nil || resp.Instance.ErrorMessage != "" {
			return nil, CreateFailedStatus, fmt.Errorf("error message: %v", resp.Instance.ErrorMessage)
		}

		if utils.StrSliceContains(targets, resp.Instance.Status) {
			return resp, RunningStatus, nil
		}

		return resp, CreatingStatus, nil
	}
}

// waitForInstanceDeleted used to check the length of instances in response
func waitForInstanceDeleted(romaConnectClient *golangsdk.ServiceClient, instanceId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := instances.CheckList(romaConnectClient, instanceId).Extract()
		if err != nil {
			return nil, ErrorStatus, err
		}

		if len(resp.Instances) == 0 {
			return resp, DeletedStatus, nil
		}
		return resp, PendingStatus, nil
	}
}
