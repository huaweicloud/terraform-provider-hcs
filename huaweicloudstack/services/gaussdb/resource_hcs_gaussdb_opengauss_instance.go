package gaussdb

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/bss/v2/orders"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/opengauss/v3/backups"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/opengauss/v3/instances"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

type HaMode string

type ConsistencyType string

const (
	HaModeDistributed HaMode = "enterprise"
	HAModeCentralized HaMode = "centralization_standard"

	ConsistencyTypeStrong   ConsistencyType = "strong"
	ConsistencyTypeEventual ConsistencyType = "eventual"
)

// GaussDB POST    /gaussdb/v3.1/{project_id}/instances
// GaussDB GET     /gaussdb/v3.1/{project_id}/{instance_id}
// GaussDB DELETE  /gaussdb/v3/{project_id}/{instance_id}
// GaussDB PUT     /gaussdb/v3/{project_id}/instances/{instance_id}/name
// GaussDB POST    /gaussdb/v3/{project_id}/instances/{instance_id}/password
// GaussDB POST    /gaussdb/v3/{project_id}/instances/{instance_id}/action
// GaussDB PUT     /gaussdb/v3/{project_id}/instances/{instance_id}/backups/policy
// GaussDB PUT     /gaussdb/v3/{project_id}/instances/{instance_id}/kms-tde/switch
func ResourceOpenGaussInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenGaussInstanceCreate,
		ReadContext:   resourceOpenGaussInstanceRead,
		UpdateContext: resourceOpenGaussInstanceUpdate,
		DeleteContext: resourceOpenGaussInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
		},
		CustomizeDiff: func(_ context.Context, d *schema.ResourceDiff, v interface{}) error {
			if d.HasChange("coordinator_num") {
				return d.SetNewComputed("private_ips")
			}
			return nil
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
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
			"ha": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(HaModeDistributed), string(HAModeCentralized),
							}, true),
							DiffSuppressFunc: utils.SuppressCaseDiffs,
						},
						"replication_mode": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"consistency": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(ConsistencyTypeStrong), string(ConsistencyTypeEventual),
							}, true),
							DiffSuppressFunc: utils.SuppressCaseDiffs,
						},
						"consistency_protocol": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"volume": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"sharding_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validation.IntBetween(1, 9),
			},
			"coordinator_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validation.IntBetween(1, 9),
			},
			"replica_num": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.IntInSlice([]int{
					2, 3,
				}),
				Default: 3,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"configuration_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "UTC+08:00",
			},
			"solution": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dorado_storage_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enable_single_float_ip": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"kms_tde_key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"kms_project_name"},
			},
			"kms_project_name": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				RequiredWith: []string{"kms_tde_key_id"},
			},
			"kms_tde_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datastore": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"backup_strategy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"keep_days": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"force_import": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),

			// Attributes
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"public_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"db_user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"switch_strategy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maintenance_window": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"management_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bms_hs_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceOpenGaussDataStore(d *schema.ResourceData) instances.DataStoreOpt {
	var db instances.DataStoreOpt

	datastoreRaw := d.Get("datastore").([]interface{})
	if len(datastoreRaw) == 1 {
		datastore := datastoreRaw[0].(map[string]interface{})
		db.Type = datastore["engine"].(string)
		db.Version = datastore["version"].(string)
	} else {
		db.Type = "GaussDB(for openGauss)"
	}
	return db
}

func resourceOpenGaussBackupStrategy(d *schema.ResourceData) *instances.BackupStrategyOpt {
	var backupOpt instances.BackupStrategyOpt

	backupStrategyRaw := d.Get("backup_strategy").([]interface{})
	if len(backupStrategyRaw) == 1 {
		strategy := backupStrategyRaw[0].(map[string]interface{})
		backupOpt.StartTime = strategy["start_time"].(string)
		backupOpt.KeepDays = strategy["keep_days"].(int)
		return &backupOpt
	}

	return nil
}

func OpenGaussInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := instances.GetInstanceByID(client, instanceID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "DELETED", nil
			}
			return nil, "", err
		}

		return v, v.Status, nil
	}
}

func buildOpenGaussInstanceCreateOpts(d *schema.ResourceData,
	config *config.HcsConfig) (instances.CreateGaussDBOpts, error) {
	createOpts := instances.CreateGaussDBOpts{
		Name:                d.Get("name").(string),
		Flavor:              d.Get("flavor").(string),
		Region:              config.GetRegion(d),
		VpcId:               d.Get("vpc_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		SecurityGroupId:     d.Get("security_group_id").(string),
		Port:                d.Get("port").(string),
		EnterpriseProjectId: config.GetEnterpriseProjectID(d),
		TimeZone:            d.Get("time_zone").(string),
		AvailabilityZone:    d.Get("availability_zone").(string),
		ConfigurationId:     d.Get("configuration_id").(string),
		OsType:              d.Get("os_type").(string),
		Solution:            d.Get("solution").(string),
		DoradoStoragePoolId: d.Get("dorado_storage_pool_id").(string),
		EnableSingleFloatIp: d.Get("enable_single_float_ip").(bool),
		ShardingNum:         d.Get("sharding_num").(int),
		CoordinatorNum:      d.Get("coordinator_num").(int),
		ReplicaNum:          d.Get("replica_num").(int),
		KmsTdeKeyId:         d.Get("kms_tde_key_id").(string),
		KmsProjectName:      d.Get("kms_project_name").(string),
		DataStore:           resourceOpenGaussDataStore(d),
		BackupStrategy:      resourceOpenGaussBackupStrategy(d),
	}

	// build HA parameter
	haRaw := d.Get("ha").([]interface{})
	log.Printf("[DEBUG] The HA structure is: %#v", haRaw)
	ha := haRaw[0].(map[string]interface{})
	mode := ha["mode"].(string)
	createOpts.Ha = &instances.HaOpt{
		Mode:                mode,
		ReplicationMode:     ha["replication_mode"].(string),
		Consistency:         ha["consistency"].(string),
		ConsistencyProtocol: ha["consistency_protocol"].(string),
	}

	// build volume
	var dn_num int = 1
	if mode == string(HaModeDistributed) {
		dn_num = d.Get("sharding_num").(int)
	}
	if mode == string(HAModeCentralized) {
		dn_num = d.Get("replica_num").(int) + 1
	}

	volumeRaw := d.Get("volume").([]interface{})
	if len(volumeRaw) > 0 {
		log.Printf("[DEBUG] The volume structure is: %#v", volumeRaw)
		volume := volumeRaw[0].(map[string]interface{})
		dn_size := volume["size"].(int)
		volume_size := dn_size * dn_num
		createOpts.Volume = instances.VolumeOpt{
			Type: volume["type"].(string),
			Size: volume_size,
		}
	}
	log.Printf("[DEBUG] The createOpts object is: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	if d.Get("charging_mode").(string) == "prePaid" {
		if err := common.ValidatePrePaidChargeInfo(d); err != nil {
			return createOpts, err
		}
		createOpts.ChargeInfo = &instances.ChargeInfo{
			ChargeMode:  "prePaid",
			PeriodType:  d.Get("period_unit").(string),
			PeriodNum:   d.Get("period").(int),
			IsAutoRenew: d.Get("auto_renew").(string),
			IsAutoPay:   "true",
		}
	}
	return createOpts, nil
}

func resourceOpenGaussInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := config.GetHcsConfig(meta)
	client, err := conf.OpenGaussV31Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GaussDB v3 client: %s ", err)
	}

	// If force_import set, try to import it instead of creating
	if common.HasFilledOpt(d, "force_import") {
		log.Printf("[DEBUG] Gaussdb opengauss instance force_import is set, try to import it instead of creating")
		listOpts := instances.ListGaussDBInstanceOpts{
			Name: d.Get("name").(string),
		}
		pages, err := instances.List(client, listOpts).AllPages()
		if err != nil {
			return diag.FromErr(err)
		}

		allInstances, err := instances.ExtractGaussDBInstances(pages)
		if err != nil {
			return diag.Errorf("unable to retrieve instances: %s", err)
		}
		if allInstances.TotalCount > 0 {
			instance := allInstances.Instances[0]
			log.Printf("[DEBUG] found existing opengauss instance %s with name %s", instance.Id, instance.Name)
			d.SetId(instance.Id)
			return resourceOpenGaussInstanceRead(ctx, d, meta)
		}
	}

	createOpts, err := buildOpenGaussInstanceCreateOpts(d, conf)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating OpenGauss instance: %s", err)
	}

	if resp.OrderId != "" {
		bssClient, err := conf.BssV2Client(conf.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, resp.OrderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, resp.OrderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(resourceId)
	} else {
		d.SetId(resp.Instance.Id)
	}

	// waiting for the instance to become ready
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"BUILD"},
		Target:                    []string{"ACTIVE", "BACKING UP"},
		Refresh:                   OpenGaussInstanceStateRefreshFunc(client, d.Id()),
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		Delay:                     20 * time.Second,
		PollInterval:              20 * time.Second,
		ContinuousTargetOccurence: 2,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to become ready: %s", d.Id(), err)
	}

	return resourceOpenGaussInstanceRead(ctx, d, meta)
}

func flattenOpenGaussDataStore(dataStore instances.DataStoreOpt) []map[string]interface{} {
	if dataStore == (instances.DataStoreOpt{}) {
		return nil
	}
	return []map[string]interface{}{
		{
			"version": dataStore.Version,
			"engine":  dataStore.Type,
		},
	}
}

func flattenOpenGaussBackupStrategy(backupStrategy instances.BackupStrategyOpt) []map[string]interface{} {
	if backupStrategy == (instances.BackupStrategyOpt{}) {
		return nil
	}
	return []map[string]interface{}{
		{
			"start_time": backupStrategy.StartTime,
			"keep_days":  backupStrategy.KeepDays,
		},
	}
}

func flattenOpenGaussVolume(volume instances.VolumeOpt, dnNum int) []map[string]interface{} {
	if volume == (instances.VolumeOpt{}) {
		return nil
	}

	return []map[string]interface{}{
		{
			"type": volume.Type,
			"size": volume.Size / dnNum,
		},
	}
}

func setOpenGaussNodesAndRelatedNumbers(d *schema.ResourceData, instance instances.GaussDBInstance,
	dnNum *int) error {
	var (
		shardingNum    = 0
		coordinatorNum = 0
	)

	nodesList := make([]map[string]interface{}, 0, 1)
	for _, raw := range instance.Nodes {
		node := map[string]interface{}{
			"id":                raw.Id,
			"name":              raw.Name,
			"status":            raw.Status,
			"role":              raw.Role,
			"availability_zone": raw.AvailabilityZone,
			"private_ip":        raw.PrivateIp,
			"public_ip":         raw.PublicIp,
			"data_ip":           raw.DataIp,
			"management_ip":     raw.ManagementIp,
			"bms_hs_ip":         raw.BmsHsIp,
		}
		nodesList = append(nodesList, node)

		if strings.Contains(raw.Name, "_gaussdbv5cn") {
			coordinatorNum += 1
		} else if strings.Contains(raw.Name, "_gaussdbv5dn") {
			shardingNum += 1
		}
	}

	if shardingNum > 0 && coordinatorNum > 0 {
		*dnNum = shardingNum / d.Get("replica_num").(int)
		return multierror.Append(nil,
			d.Set("nodes", nodesList),
			d.Set("sharding_num", dnNum),
			d.Set("coordinator_num", coordinatorNum),
		).ErrorOrNil()
	} else {
		// If the HA mode is centralized, the HA structure of API response is nil.
		*dnNum = instance.ReplicaNum + 1
		return multierror.Append(nil,
			d.Set("nodes", nodesList),
			d.Set("replica_num", instance.ReplicaNum),
		).ErrorOrNil()
	}
}

func setOpenGaussPrivateIpsAndEndpoints(d *schema.ResourceData, privateIps []string, port int) error {
	if len(privateIps) < 1 {
		return nil
	}

	privateIp := privateIps[0]
	ip_list := strings.Split(privateIp, "/")
	endpoints := []string{}
	for i := 0; i < len(ip_list); i++ {
		ip_list[i] = strings.Trim(ip_list[i], " ")
		endpoint := fmt.Sprintf("%s:%d", ip_list[i], port)
		endpoints = append(endpoints, endpoint)
	}

	return multierror.Append(nil,
		d.Set("private_ips", ip_list),
		d.Set("endpoints", endpoints),
	).ErrorOrNil()
}

func resourceOpenGaussInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := config.GetHcsConfig(meta)
	region := conf.GetRegion(d)
	client, err := conf.OpenGaussV31Client(region)
	if err != nil {
		return diag.Errorf("error creating GaussDB v3.1 client: %s ", err)
	}

	instanceID := d.Id()
	instance, err := instances.GetInstanceByID(client, instanceID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "OpenGauss instance")
	}
	if instance.Id == "" {
		d.SetId("")
		return nil
	}

	var dnNum int = 1
	log.Printf("[DEBUG] Retrieved instance (%s): %#v", instanceID, instance)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", instance.Name),
		d.Set("status", instance.Status),
		d.Set("type", instance.Type),
		d.Set("vpc_id", instance.VpcId),
		d.Set("subnet_id", instance.SubnetId),
		d.Set("security_group_id", instance.SecurityGroupId),
		d.Set("db_user_name", instance.DbUserName),
		d.Set("time_zone", instance.TimeZone),
		d.Set("flavor", instance.FlavorRef),
		d.Set("port", strconv.Itoa(instance.Port)),
		d.Set("switch_strategy", instance.SwitchStrategy),
		d.Set("maintenance_window", instance.MaintenanceWindow),
		d.Set("public_ips", instance.PublicIps),
		d.Set("charging_mode", instance.ChargeInfo.ChargeMode),
		d.Set("datastore", flattenOpenGaussDataStore(instance.DataStore)),
		d.Set("backup_strategy", flattenOpenGaussBackupStrategy(instance.BackupStrategy)),
		setOpenGaussNodesAndRelatedNumbers(d, instance, &dnNum),
		d.Set("volume", flattenOpenGaussVolume(instance.Volume, dnNum)),
		setOpenGaussPrivateIpsAndEndpoints(d, instance.PrivateIps, instance.Port),
		d.Set("kms_tde_key_id", instance.KmsTdeKeyId),
		d.Set("kms_project_name", instance.KmsProjectName),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting OpenGauss instance fields: %s", mErr.ErrorOrNil())
	}
	return nil
}

func expandOpenGaussShardingNumber(ctx context.Context, config *config.HcsConfig, client *golangsdk.ServiceClient,
	d *schema.ResourceData) error {
	old, newnum := d.GetChange("sharding_num")
	if newnum.(int) < old.(int) {
		return fmt.Errorf("error expanding shard for instance: new num must be larger than the old one.")
	}
	expandSize := newnum.(int) - old.(int)
	opts := instances.UpdateOpts{
		ExpandCluster: &instances.UpdateClusterOpts{
			Shard: &instances.Shard{
				Count: expandSize,
			},
		},
		IsAutoPay: "true",
	}
	log.Printf("[DEBUG] The updateOpts object of sharding number is: %#v", opts)
	return updateVolumeAndRelatedHaNumbers(ctx, config, client, d, opts)
}

func expandOpenGaussCoordinatorNumber(ctx context.Context, config *config.HcsConfig, client *golangsdk.ServiceClient,
	d *schema.ResourceData) error {
	old, newnum := d.GetChange("coordinator_num")
	if newnum.(int) < old.(int) {
		return fmt.Errorf("error expanding coordinator for instance: new number must be larger than the old one.")
	}
	expandSize := newnum.(int) - old.(int)

	var coordinators []instances.Coordinator
	azlist := strings.Split(d.Get("availability_zone").(string), ",")
	for i := 0; i < expandSize; i++ {
		coordinator := instances.Coordinator{
			AzCode: azlist[0],
		}
		coordinators = append(coordinators, coordinator)
	}
	opts := instances.UpdateOpts{
		ExpandCluster: &instances.UpdateClusterOpts{
			Coordinators: coordinators,
		},
		IsAutoPay: "true",
	}
	log.Printf("[DEBUG] The updateOpts object of coordinator number is: %#v", opts)
	return updateVolumeAndRelatedHaNumbers(ctx, config, client, d, opts)
}

func updateOpenGaussVolumeSize(ctx context.Context, config *config.HcsConfig, client *golangsdk.ServiceClient,
	d *schema.ResourceData) error {
	volumeRaw := d.Get("volume").([]interface{})
	dnSize := volumeRaw[0].(map[string]interface{})["size"].(int)
	dnNum := 1
	if d.Get("ha.0.mode").(string) == string(HaModeDistributed) {
		dnNum = d.Get("sharding_num").(int)
	}
	if d.Get("ha.0.mode").(string) == string(HAModeCentralized) {
		dnNum = d.Get("replica_num").(int) + 1
	}
	opts := instances.UpdateOpts{
		EnlargeVolume: &instances.UpdateVolumeOpts{
			Size: dnSize * dnNum,
		},
		IsAutoPay: "true",
	}
	log.Printf("[DEBUG] The updateOpts object of volume size is: %#v", opts)
	return updateVolumeAndRelatedHaNumbers(ctx, config, client, d, opts)
}

func updateVolumeAndRelatedHaNumbers(ctx context.Context, config *config.HcsConfig, client *golangsdk.ServiceClient,
	d *schema.ResourceData, opts instances.UpdateOpts) error {
	instanceId := d.Id()
	resp, err := instances.Update(client, instanceId, opts)
	if err != nil {
		return fmt.Errorf("error updating instance (%s): %s", instanceId, err)
	}
	if resp.OrderId != "" {
		bssClient, err := config.BssV2Client(config.GetRegion(d))
		if err != nil {
			return fmt.Errorf("error creating BSS v2 client: %s", err)
		}
		if err := orders.WaitForOrderSuccess(bssClient, int(d.Timeout(schema.TimeoutUpdate)/time.Second), resp.OrderId); err != nil {
			return err
		}
	}
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"MODIFYING", "EXPANDING", "BACKING UP"},
		Target:                    []string{"ACTIVE"},
		Refresh:                   OpenGaussInstanceStateRefreshFunc(client, instanceId),
		Timeout:                   d.Timeout(schema.TimeoutUpdate),
		Delay:                     20 * time.Second,
		PollInterval:              20 * time.Second,
		ContinuousTargetOccurence: 2,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for instance (%s) status to active: %s ", instanceId, err)
	}

	return nil
}

func resourceOpenGaussInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := config.GetHcsConfig(meta)
	region := conf.GetRegion(d)
	client, err := conf.OpenGaussV3Client(region)
	if err != nil {
		return diag.Errorf("error creating GaussDB v3 client: %s ", err)
	}

	log.Printf("[DEBUG] Updating OpenGaussDB instances %s", d.Id())
	instanceId := d.Id()

	if d.HasChange("name") {
		renameOpts := instances.RenameOpts{
			Name: d.Get("name").(string),
		}
		_, err = instances.Rename(client, renameOpts, instanceId).Extract()
		if err != nil {
			return diag.Errorf("error updating name for instance (%s): %s", instanceId, err)
		}
	}

	if d.HasChange("password") {
		restorePasswordOpts := instances.RestorePasswordOpts{
			Password: d.Get("password").(string),
		}
		r := golangsdk.ErrResult{}
		r.Result = instances.RestorePassword(client, restorePasswordOpts, instanceId)
		if r.ExtractErr() != nil {
			return diag.Errorf("error updating password for instance (%s): %s ", instanceId, r.Err)
		}
	}

	if d.HasChange("sharding_num") {
		if err := expandOpenGaussShardingNumber(ctx, conf, client, d); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("coordinator_num") {
		if err := expandOpenGaussCoordinatorNumber(ctx, conf, client, d); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("volume") {
		if err := updateOpenGaussVolumeSize(ctx, conf, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("backup_strategy") {
		backupRaw := d.Get("backup_strategy").([]interface{})
		rawMap := backupRaw[0].(map[string]interface{})
		keep_days := rawMap["keep_days"].(int)

		updateOpts := backups.UpdateOpts{
			KeepDays:           &keep_days,
			StartTime:          rawMap["start_time"].(string),
			Period:             "1,2,3,4,5,6,7", // Fixed to "1,2,3,4,5,6,7"
			DifferentialPeriod: "30",            // Fixed to "30"
		}

		log.Printf("[DEBUG] The updateOpts object of backup_strategy parameter is: %#v", updateOpts)
		err = backups.Update(client, d.Id(), updateOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error updating backup_strategy: %s", err)
		}
	}

	if d.HasChanges("kms_tde_key_id", "kms_project_name", "kms_tde_status") {
		updateKmsOpts := instances.UpdateKmsOpts{
			KmsTdeKeyId:    d.Get("kms_tde_key_id").(string),
			KmsProjectName: d.Get("kms_project_name").(string),
			KmsTdeStatus:   d.Get("kms_tde_status").(string),
		}
		log.Printf("[DEBUG] The update object of KMS is: %#v", updateKmsOpts)
		_, err = instances.UpdateKms(client, updateKmsOpts, d.Id()).Extract()
		if err != nil {
			return diag.Errorf("error updating KMS information for instance (%s): %s", d.Id(), err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := conf.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), d.Id()); err != nil {
			return diag.Errorf("error updating the auto-renew of the instance (%s): %s", d.Id(), err)
		}
	}

	return resourceOpenGaussInstanceRead(ctx, d, meta)
}

func resourceOpenGaussInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := config.GetHcsConfig(meta)
	client, err := conf.OpenGaussV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GaussDB v3 client: %s ", err)
	}

	instanceId := d.Id()
	if v, ok := d.GetOk("charging_mode"); ok && v.(string) == "prePaid" {
		if err := common.UnsubscribePrePaidResource(d, conf, []string{instanceId}); err != nil {
			return diag.Errorf("error unsubscribe OpenGauss instance: %s", err)
		}
	} else {
		result := instances.Delete(client, instanceId)
		if result.Err != nil {
			return common.CheckDeletedDiag(d, result.Err, "OpenGauss instance")
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:        []string{"ACTIVE", "BACKING UP", "FAILED"},
		Target:         []string{"DELETED"},
		Refresh:        OpenGaussInstanceStateRefreshFunc(client, instanceId),
		Timeout:        d.Timeout(schema.TimeoutDelete),
		Delay:          60 * time.Second,
		MinTimeout:     20 * time.Second,
		NotFoundChecks: 2,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for instance (%s) to be deleted: %s", instanceId, err)
	}
	log.Printf("[DEBUG] Instance deleted successfully %s", instanceId)
	return nil
}
