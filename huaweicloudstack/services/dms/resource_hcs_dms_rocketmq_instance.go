package dms

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type dmsError struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

// GetResponse response
type GetResponse struct {
	RegionID       string          `json:"regionId"`
	AvailableZones []AvailableZone `json:"available_zones"`
}

// AvailableZone for dms
type AvailableZone struct {
	ID                   string `json:"id"`
	Code                 string `json:"code"`
	Name                 string `json:"name"`
	Port                 string `json:"port"`
	ResourceAvailability string `json:"resource_availability"`
	SoldOut              bool   `json:"soldOut"`
	DefaultAz            bool   `json:"default_az"`
	RemainTime           uint64 `json:"remain_time"`
	Ipv6Enable           bool   `json:"ipv6_enable"`
}

func ResourceDmsRocketMQInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRocketMQInstanceCreate,
		UpdateContext: resourceDmsRocketMQInstanceUpdate,
		ReadContext:   resourceDmsRocketMQInstanceRead,
		DeleteContext: resourceDmsRocketMQInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Update: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the DMS RocketMQ instance`,
			},
			"engine_version": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the version of the RocketMQ engine.`,
			},
			"storage_space": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the message storage capacity, Unit: GB.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of a VPC`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of a subnet`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of a security group`,
			},
			"availability_zones": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				ForceNew:    true,
				Set:         schema.HashString,
				Description: `Specifies the list of availability zone names`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies a product ID`,
			},
			"storage_spec_code": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the storage I/O specification`,
			},
			"broker_num": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the broker numbers.(HCS is required, HC is optional)`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the description of the DMS RocketMQ instance.`,
			},
			"ssl_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies whether the RocketMQ SASL_SSL is enabled.`,
			},
			"ipv6_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies whether to support IPv6`,
			},
			"enable_publicip": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies whether to enable public access.`,
			},
			"publicip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the EIP bound to the instance.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the enterprise project id of the instance.`,
			},
			"enable_acl": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether access control is enabled.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the DMS RocketMQ instance.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the DMS RocketMQ instance type. Value: cluster.`,
			},
			"specification": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Indicates the instance specification. For a cluster DMS RocketMQ instance, VM specifications
  and the number of nodes are returned.`,
			},
			"maintain_begin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time at which the maintenance window starts. The format is HH:mm:ss.`,
			},
			"maintain_end": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time at which the maintenance window ends. The format is HH:mm:ss.`,
			},
			"used_storage_space": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the used message storage space. Unit: GB.`,
			},
			"publicip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the public IP address.`,
			},
			"node_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the node quantity.`,
			},
			"new_spec_billing_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether billing based on new specifications is enabled.`,
			},
			"namesrv_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the metadata address.`,
			},
			"broker_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the service data address.`,
			},
			"public_namesrv_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the public network metadata address.`,
			},
			"public_broker_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the public network service data address.`,
			},
			"resource_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the resource specifications.`,
			},
			"retention_policy": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				Deprecated: "Use 'enable_acl' instead",
			},
			"tags": common.TagsSchema(),
			"cross_vpc_accesses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"advertised_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// Typo, it is only kept in the code, will not be shown in the docs.
						"lisenter_ip": {
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "typo in lisenter_ip, please use \"listener_ip\" instead.",
						},
					},
				},
			},
		},
	}
}

func resourceDmsRocketMQInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createRocketmqInstance: create DMS rocketmq instance
	var (
		createRocketmqInstanceHttpUrl = "v2/{project_id}/instances"
		createRocketmqInstanceProduct = "dmsv2"
	)
	createRocketmqInstanceClient, err := cfg.NewServiceClient(createRocketmqInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQInstance Client: %s", err)
	}

	createRocketmqInstancePath := createRocketmqInstanceClient.Endpoint + createRocketmqInstanceHttpUrl
	createRocketmqInstancePath = strings.ReplaceAll(createRocketmqInstancePath, "{project_id}",
		createRocketmqInstanceClient.ProjectID)

	createRocketmqInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	var availableZones []string
	// convert the codes of the availability zone into ids
	azCodes := d.Get("availability_zones").(*schema.Set)
	availableZones, err = getAvailableZoneIDByCode(cfg, region, azCodes.List())
	if err != nil {
		return diag.FromErr(err)
	}
	createRocketmqInstanceOpt.JSONBody = utils.RemoveNil(buildCreateRocketmqInstanceBodyParams(d, cfg, availableZones))
	createRocketmqInstanceResp, err := createRocketmqInstanceClient.Request("POST", createRocketmqInstancePath,
		&createRocketmqInstanceOpt)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQInstance: %s", err)
	}
	createRocketmqInstanceRespBody, err := utils.FlattenResponse(createRocketmqInstanceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("instance_id", createRocketmqInstanceRespBody)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQInstance: ID is not found in API response")
	}

	var delayTime time.Duration = 300
	if chargingMode, ok := d.GetOk("charging_mode"); ok && chargingMode == "prePaid" {
		err = waitForRocketMQOrderComplete(ctx, d, cfg, createRocketmqInstanceClient, id.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		delayTime = 5
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"RUNNING"},
		Refresh:      rocketmqInstanceStateRefreshFunc(createRocketmqInstanceClient, id.(string)),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        delayTime * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to create: %s", id.(string), err)
	}

	d.SetId(id.(string))

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		tagList := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(createRocketmqInstanceClient, "rocketmq", id.(string), tagList).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of RocketMQ %s: %s", id.(string), tagErr)
		}
	}

	return resourceDmsRocketMQInstanceRead(ctx, d, meta)
}

func waitForRocketMQOrderComplete(ctx context.Context, d *schema.ResourceData, cfg *config.Config,
	client *golangsdk.ServiceClient, instanceID string) error {
	region := cfg.GetRegion(d)
	orderId, err := getRocketMQInstanceOrderId(ctx, d, client, instanceID)
	if err != nil {
		return err
	}

	if orderId == "" {
		log.Printf("[WARN] error get order id by instance ID: %s", instanceID)
		return nil
	}

	bssClient, err := cfg.BssV2Client(region)
	if err != nil {
		return fmt.Errorf("error creating BSS v2 client: %s", err)
	}
	// wait for order success
	err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	_, err = common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error waiting for RocketMQ order resource %s complete: %s", orderId, err)
	}
	return nil
}

func getRocketMQInstanceOrderId(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string) (string, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"EMPTY"},
		Target:       []string{"CREATING"},
		Refresh:      rocketMQInstanceCreatingFunc(client, instanceID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        500 * time.Millisecond,
		PollInterval: 500 * time.Millisecond,
	}
	orderId, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return "", fmt.Errorf("error waiting for RocketMQ instance (%s) to creating: %s", instanceID, err)
	}
	return orderId.(string), nil
}

func rocketMQInstanceCreatingFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			getRocketmqInstanceHttpUrl = "v2/{project_id}/instances/{instance_id}"
		)

		getRocketmqInstancePath := client.Endpoint + getRocketmqInstanceHttpUrl
		getRocketmqInstancePath = strings.ReplaceAll(getRocketmqInstancePath, "{project_id}", client.ProjectID)
		getRocketmqInstancePath = strings.ReplaceAll(getRocketmqInstancePath, "{instance_id}", fmt.Sprintf("%v", instanceID))

		getRocketmqInstanceOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		getRocketmqInstanceResp, err := client.Request("GET", getRocketmqInstancePath, &getRocketmqInstanceOpt)

		if err != nil {
			if errCode, ok := err.(golangsdk.ErrDefault404); ok {
				var rocketMQError dmsError
				err = json.Unmarshal(errCode.Body, &rocketMQError)
				if err != nil {
					return nil, "", fmt.Errorf("error get DmsRocketMQInstance: error format error: %s", err)
				}
				if rocketMQError.ErrorCode == "DMS.00404022" {
					return getRocketmqInstanceResp, "EMPTY", nil
				}
			}
			return nil, "", fmt.Errorf("error retrieving DmsRocketMQInstance: %s", err)
		}

		res, err := utils.FlattenResponse(getRocketmqInstanceResp)
		if err != nil {
			return nil, "", err
		}
		orderID := utils.PathSearch("order_id", res, "")
		return orderID, "CREATING", nil
	}
}

func buildCreateRocketmqInstanceBodyParams(d *schema.ResourceData, cfg *config.Config,
	availableZones []string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                  utils.ValueIngoreEmpty(d.Get("name")),
		"enable_acl":            utils.ValueIngoreEmpty(d.Get("enable_acl")),
		"description":           utils.ValueIngoreEmpty(d.Get("description")),
		"engine":                "reliability",
		"engine_version":        utils.ValueIngoreEmpty(d.Get("engine_version")),
		"storage_space":         utils.ValueIngoreEmpty(d.Get("storage_space")),
		"vpc_id":                utils.ValueIngoreEmpty(d.Get("vpc_id")),
		"subnet_id":             utils.ValueIngoreEmpty(d.Get("subnet_id")),
		"security_group_id":     utils.ValueIngoreEmpty(d.Get("security_group_id")),
		"available_zones":       availableZones,
		"product_id":            utils.ValueIngoreEmpty(d.Get("flavor_id")),
		"ssl_enable":            utils.ValueIngoreEmpty(d.Get("ssl_enable")),
		"storage_spec_code":     utils.ValueIngoreEmpty(d.Get("storage_spec_code")),
		"ipv6_enable":           utils.ValueIngoreEmpty(d.Get("ipv6_enable")),
		"enable_publicip":       utils.ValueIngoreEmpty(d.Get("enable_publicip")),
		"publicip_id":           utils.ValueIngoreEmpty(d.Get("publicip_id")),
		"broker_num":            utils.ValueIngoreEmpty(d.Get("broker_num")),
		"enterprise_project_id": utils.ValueIngoreEmpty(common.GetEnterpriseProjectID(d, cfg)),
	}
	if chargingMode, ok := d.GetOk("charging_mode"); ok && chargingMode == "prePaid" {
		bodyParams["bss_param"] = buildCreateRocketmqInstanceBodyBssParams(d)
	}
	return bodyParams
}

func buildCreateRocketmqInstanceBodyBssParams(d *schema.ResourceData) map[string]interface{} {
	var autoRenew bool
	if d.Get("auto_renew").(string) == "true" {
		autoRenew = true
	}
	isAutoPay := true
	bodyParams := map[string]interface{}{
		"charging_mode": utils.ValueIngoreEmpty(d.Get("charging_mode")),
		"period_type":   utils.ValueIngoreEmpty(d.Get("period_unit")),
		"period_num":    utils.ValueIngoreEmpty(d.Get("period")),
		"is_auto_renew": &autoRenew,
		"is_auto_pay":   &isAutoPay,
	}
	return bodyParams
}

func resourceDmsRocketMQInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateRocketmqInstanceHasChanges := []string{
		"name",
		"description",
		"security_group_id",
		"retention_policy",
		"enable_acl",
		"cross_vpc_accesses",
		"auto_renew",
		"tags",
	}

	if d.HasChanges(updateRocketmqInstanceHasChanges...) {
		// updateRocketmqInstance: update DMS rocketmq instance
		var (
			updateRocketmqInstanceHttpUrl = "v2/{project_id}/instances/{instance_id}"
			updateRocketmqInstanceProduct = "dmsv2"
		)
		updateRocketmqInstanceClient, err := cfg.NewServiceClient(updateRocketmqInstanceProduct, region)
		if err != nil {
			return diag.Errorf("error creating DmsRocketMQInstance Client: %s", err)
		}

		updateRocketmqInstancePath := updateRocketmqInstanceClient.Endpoint + updateRocketmqInstanceHttpUrl
		updateRocketmqInstancePath = strings.ReplaceAll(updateRocketmqInstancePath, "{project_id}",
			updateRocketmqInstanceClient.ProjectID)
		updateRocketmqInstancePath = strings.ReplaceAll(updateRocketmqInstancePath, "{instance_id}", fmt.Sprintf("%v", d.Id()))

		updateRocketmqInstanceOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				204,
			},
		}
		updateRocketmqInstanceOpt.JSONBody = utils.RemoveNil(buildUpdateRocketmqInstanceBodyParams(d))
		_, err = updateRocketmqInstanceClient.Request("PUT", updateRocketmqInstancePath, &updateRocketmqInstanceOpt)
		if err != nil {
			return diag.Errorf("error updating DmsRocketMQInstance: %s", err)
		}

		if d.HasChange("auto_renew") {
			bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
			if err != nil {
				return diag.Errorf("error creating BSS V2 client: %s", err)
			}
			if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), d.Id()); err != nil {
				return diag.Errorf("error updating the auto-renew of the RocketMQ instance (%s): %s", d.Id(), err)
			}
		}
		// update tags
		if d.HasChange("tags") {
			tagErr := utils.UpdateResourceTags(updateRocketmqInstanceClient, d, "rocketmq", d.Id())
			if tagErr != nil {
				return diag.Errorf("error updating tags of RocketMQ:%s, err:%s", d.Id(), tagErr)
			}
		}
	}

	return resourceDmsRocketMQInstanceRead(ctx, d, meta)
}

func buildUpdateRocketmqInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"description":       utils.ValueIngoreEmpty(d.Get("description")),
		"security_group_id": utils.ValueIngoreEmpty(d.Get("security_group_id")),
	}

	if d.HasChange("enable_acl") {
		bodyParams["enable_acl"] = utils.ValueIngoreEmpty(d.Get("enable_acl"))
	} else if d.HasChange("retention_policy") {
		bodyParams["enable_acl"] = utils.ValueIngoreEmpty(d.Get("retention_policy"))
	}

	if d.HasChange("name") {
		bodyParams["name"] = utils.ValueIngoreEmpty(d.Get("name"))
	}
	return bodyParams
}

func resourceDmsRocketMQInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getRocketmqInstance: Query DMS rocketmq instance
	var (
		getRocketmqInstanceHttpUrl = "v2/{project_id}/instances/{instance_id}"
		getRocketmqInstanceProduct = "dmsv2"
	)
	getRocketmqInstanceClient, err := cfg.NewServiceClient(getRocketmqInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQInstance Client: %s", err)
	}

	getRocketmqInstancePath := getRocketmqInstanceClient.Endpoint + getRocketmqInstanceHttpUrl
	getRocketmqInstancePath = strings.ReplaceAll(getRocketmqInstancePath, "{project_id}",
		getRocketmqInstanceClient.ProjectID)
	getRocketmqInstancePath = strings.ReplaceAll(getRocketmqInstancePath, "{instance_id}", fmt.Sprintf("%v", d.Id()))

	getRocketmqInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRocketmqInstanceResp, err := getRocketmqInstanceClient.Request("GET", getRocketmqInstancePath, &getRocketmqInstanceOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DmsRocketMQInstance")
	}

	getRocketmqInstanceRespBody, err := utils.FlattenResponse(getRocketmqInstanceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// convert the ids of the availability zone into codes
	var availableZoneCodes []string
	availableZoneIDs := utils.PathSearch("available_zones", getRocketmqInstanceRespBody, nil)
	if availableZoneIDs != nil {
		azIDs := make([]string, 0)
		for _, v := range availableZoneIDs.([]interface{}) {
			azIDs = append(azIDs, v.(string))
		}
		availableZoneCodes, err = getAvailableZoneCodeByID(cfg, region, azIDs)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	crossVpcInfo := utils.PathSearch("cross_vpc_info", getRocketmqInstanceRespBody, nil)
	var crossVpcAccess []map[string]interface{}
	if crossVpcInfo != nil {
		crossVpcAccess, err = flattenCrossVpcInfo(crossVpcInfo.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	var chargingMode = "postPaid"
	if utils.PathSearch("charging_mode", getRocketmqInstanceRespBody, 1).(float64) == 0 {
		chargingMode = "prePaid"
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getRocketmqInstanceRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRocketmqInstanceRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRocketmqInstanceRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getRocketmqInstanceRespBody, nil)),
		d.Set("specification", utils.PathSearch("specification", getRocketmqInstanceRespBody, nil)),
		d.Set("engine_version", utils.PathSearch("engine_version", getRocketmqInstanceRespBody, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", getRocketmqInstanceRespBody, nil)),
		d.Set("flavor_id", utils.PathSearch("product_id", getRocketmqInstanceRespBody, nil)),
		d.Set("security_group_id", utils.PathSearch("security_group_id", getRocketmqInstanceRespBody, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", getRocketmqInstanceRespBody, nil)),
		d.Set("availability_zones", availableZoneCodes),
		d.Set("maintain_begin", utils.PathSearch("maintain_begin", getRocketmqInstanceRespBody, nil)),
		d.Set("maintain_end", utils.PathSearch("maintain_end", getRocketmqInstanceRespBody, nil)),
		d.Set("storage_space", utils.PathSearch("total_storage_space", getRocketmqInstanceRespBody, nil)),
		d.Set("used_storage_space", utils.PathSearch("used_storage_space", getRocketmqInstanceRespBody, nil)),
		d.Set("enable_publicip", utils.PathSearch("enable_publicip", getRocketmqInstanceRespBody, nil)),
		d.Set("publicip_id", utils.PathSearch("publicip_id", getRocketmqInstanceRespBody, nil)),
		d.Set("publicip_address", utils.PathSearch("publicip_address", getRocketmqInstanceRespBody, nil)),
		d.Set("ssl_enable", utils.PathSearch("ssl_enable", getRocketmqInstanceRespBody, nil)),
		d.Set("storage_spec_code", utils.PathSearch("storage_spec_code", getRocketmqInstanceRespBody, nil)),
		d.Set("ipv6_enable", utils.PathSearch("ipv6_enable", getRocketmqInstanceRespBody, nil)),
		d.Set("node_num", utils.PathSearch("node_num", getRocketmqInstanceRespBody, nil)),
		d.Set("new_spec_billing_enable", utils.PathSearch("new_spec_billing_enable", getRocketmqInstanceRespBody, nil)),
		d.Set("enable_acl", utils.PathSearch("enable_acl", getRocketmqInstanceRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", getRocketmqInstanceRespBody, nil)),
		d.Set("broker_num", utils.PathSearch("broker_num", getRocketmqInstanceRespBody, nil)),
		d.Set("namesrv_address", utils.PathSearch("namesrv_address", getRocketmqInstanceRespBody, nil)),
		d.Set("broker_address", utils.PathSearch("broker_address", getRocketmqInstanceRespBody, nil)),
		d.Set("public_namesrv_address", utils.PathSearch("public_namesrv_address", getRocketmqInstanceRespBody, nil)),
		d.Set("public_broker_address", utils.PathSearch("public_broker_address", getRocketmqInstanceRespBody, nil)),
		d.Set("resource_spec_code", utils.PathSearch("resource_spec_code", getRocketmqInstanceRespBody, nil)),
		d.Set("cross_vpc_accesses", crossVpcAccess),
		d.Set("charging_mode", chargingMode),
	)

	// fetch tags
	if resourceTags, err := tags.Get(getRocketmqInstanceClient, "rocketmq", d.Id()).Extract(); err == nil {
		tagMap := utils.TagsToMap(resourceTags.Tags)
		mErr = multierror.Append(mErr, d.Set("tags", tagMap))
	} else {
		fmt.Printf("[WARN] fetching tags of RocketMQ failed: %s", err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDmsRocketMQInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteRocketmqInstance: Delete DMS rocketmq instance
	var (
		deleteRocketmqInstanceHttpUrl = "v2/{project_id}/instances/{instance_id}"
		deleteRocketmqInstanceProduct = "dmsv2"
	)
	deleteRocketmqInstanceClient, err := cfg.NewServiceClient(deleteRocketmqInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQInstance Client: %s", err)
	}

	if d.Get("charging_mode") == "prePaid" {
		if err := common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()}); err != nil {
			return diag.Errorf("error unsubscribe RocketMQ instance: %s", err)
		}
	} else {
		deleteRocketmqInstancePath := deleteRocketmqInstanceClient.Endpoint + deleteRocketmqInstanceHttpUrl
		deleteRocketmqInstancePath = strings.ReplaceAll(deleteRocketmqInstancePath, "{project_id}",
			deleteRocketmqInstanceClient.ProjectID)
		deleteRocketmqInstancePath = strings.ReplaceAll(deleteRocketmqInstancePath, "{instance_id}",
			fmt.Sprintf("%v", d.Id()))

		deleteRocketmqInstanceOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				204,
			},
		}
		_, err = deleteRocketmqInstanceClient.Request("DELETE", deleteRocketmqInstancePath, &deleteRocketmqInstanceOpt)
		if err != nil {
			return diag.Errorf("error deleting DmsRocketMQInstance: %s", err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"DELETING", "RUNNING", "ERROR"},
		Target:       []string{"DELETED"},
		Refresh:      rocketmqInstanceStateRefreshFunc(deleteRocketmqInstanceClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        90 * time.Second,
		PollInterval: 15 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to delete: %s", d.Id(), err)
	}

	d.SetId("")

	return nil
}

func rocketmqInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getRocketmqInstancePath := client.Endpoint + "v2/{project_id}/instances/{instance_id}"
		getRocketmqInstancePath = strings.ReplaceAll(getRocketmqInstancePath, "{project_id}", client.ProjectID)
		getRocketmqInstancePath = strings.ReplaceAll(getRocketmqInstancePath, "{instance_id}", fmt.Sprintf("%v", instanceID))
		getRocketmqInstanceOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		v, err := client.Request("GET", getRocketmqInstancePath, &getRocketmqInstanceOpt)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "DELETED", nil
			}
			return nil, "", err
		}
		respBody, err := utils.FlattenResponse(v)
		if err != nil {
			return nil, "", err
		}
		status := utils.PathSearch("status", respBody, "").(string)
		return respBody, status, nil
	}
}

func flattenCrossVpcInfo(str string) (result []map[string]interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Recover panic when flattening Cross-VPC structure: %#v \nCrossVpcInfo: %s", r, str)
			err = fmt.Errorf("faield to flattening Cross-VPC structure: %#v", r)
		}
	}()

	return unmarshalFlattenCrossVpcInfo(str)
}

func unmarshalFlattenCrossVpcInfo(crossVpcInfoStr string) ([]map[string]interface{}, error) {
	if crossVpcInfoStr == "" {
		return nil, nil
	}

	crossVpcInfos := make(map[string]interface{})
	err := json.Unmarshal([]byte(crossVpcInfoStr), &crossVpcInfos)
	if err != nil {
		return nil, fmt.Errorf("failed to Unmarshal CrossVpcInfo, crossVpcInfo: %s, error: %s", crossVpcInfoStr, err)
	}

	ipArr := make([]string, 0, len(crossVpcInfos))
	for ip := range crossVpcInfos {
		ipArr = append(ipArr, ip)
	}
	sort.Strings(ipArr) // Sort by listeners IP.

	result := make([]map[string]interface{}, len(crossVpcInfos))
	for i, ip := range ipArr {
		crossVpcInfo := crossVpcInfos[ip].(map[string]interface{})
		result[i] = map[string]interface{}{
			"listener_ip":   ip,
			"lisenter_ip":   ip,
			"advertised_ip": crossVpcInfo["advertised_ip"],
			"port":          crossVpcInfo["port"],
			"port_id":       crossVpcInfo["port_id"],
		}
	}
	return result, nil
}

func getAvailableZoneIDByCode(config *config.Config, region string, azCodes []interface{}) ([]string, error) {
	if len(azCodes) == 0 {
		return nil, fmt.Errorf(`arguments "azCodes" is required`)
	}

	availableZones, err := getAvailableZones(config, region)
	if err != nil {
		return nil, err
	}

	codeIDMapping := make(map[string]string)
	for _, v := range availableZones {
		codeIDMapping[v.Code] = v.ID
	}

	azIDs := make([]string, 0, len(azCodes))
	for _, code := range azCodes {
		if id, ok := codeIDMapping[code.(string)]; ok {
			azIDs = append(azIDs, id)
		}
	}
	log.Printf("[DEBUG] DMS converts the AZ codes to AZ IDs: \n%#v => \n%#v", azCodes, azIDs)
	return azIDs, nil
}

func getAvailableZoneCodeByID(config *config.Config, region string, azIDs []string) ([]string, error) {
	if len(azIDs) == 0 {
		return nil, fmt.Errorf(`arguments "azIDs" is required`)
	}

	availableZones, err := getAvailableZones(config, region)
	if err != nil {
		return nil, err
	}

	idCodeMapping := make(map[string]string)
	for _, v := range availableZones {
		idCodeMapping[v.ID] = v.Code
	}

	azCodes := make([]string, 0, len(azIDs))
	for _, id := range azIDs {
		if code, ok := idCodeMapping[id]; ok {
			azCodes = append(azCodes, code)
		}
	}
	log.Printf("[DEBUG] DMS converts the AZ IDs to AZ codes: \n%#v => \n%#v", azIDs, azCodes)
	return azCodes, nil
}

func getAvailableZones(cfg *config.Config, region string) ([]AvailableZone, error) {
	client, err := cfg.DmsV2Client(region)
	if err != nil {
		return nil, fmt.Errorf("error initializing DMS(v2) client: %s", err)
	}

	r, err := azGet(client)
	if err != nil {
		return nil, fmt.Errorf("error querying available Zones: %s", err)
	}

	return r.AvailableZones, nil
}

// Get available zones
func azGet(client *golangsdk.ServiceClient) (*GetResponse, error) {
	var rst golangsdk.Result
	_, err := client.Get(getURL(client), &rst.Body, nil)
	if err == nil {
		var r GetResponse
		err = rst.ExtractInto(&r)
		return &r, err
	}
	return nil, err
}

// getURL will build the get url of get function
func getURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("available-zones")
}
