// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DNS
// ---------------------------------------------------------------

package dns

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/dns/v2/zones"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func ResourceDNSRecordset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSRecordsetCreate,
		UpdateContext: resourceDNSRecordsetUpdate,
		ReadContext:   resourceDNSRecordsetRead,
		DeleteContext: resourceDNSRecordsetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the zone ID.`,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Specifies the name of the record set. The name suffixed with a zone name, which is a
complete host name ended with a dot.`,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"A", "AAAA", "MX", "CNAME", "TXT", "NS", "SRV",
				}, false),
				Description: `Specifies the type of the record set.`,
			},
			"records": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
				Required:    true,
				Description: `Specifies an array of DNS records. The value rules vary depending on the record set type.`,
			},
			"ttl": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      300,
				ValidateFunc: validation.IntBetween(1, 2147483647),
				Description:  `Specifies the time to live (TTL) of the record set (in seconds).`,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ENABLE",
				ValidateFunc: validation.StringInSlice([]string{"ENABLE", "DISABLE"}, false),
				Description:  `Specifies the status of the record set.`,
			},
			"tags": common.TagsSchema(),
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the description of the record set.`,
			},
			"zone_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The zone name of the record set.`,
			},
		},
	}
}

type WaitForConfig struct {
	ZoneID      string
	RecordsetID string
	ZoneType    string
	Timeout     time.Duration
}

func resourceDNSRecordsetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	createDNSRecordsetClient, err := cfg.NewServiceClient("dns_region", region)
	if err != nil {
		return diag.Errorf("error creating DNS Client: %s", err)
	}

	zoneType, err := getDNSZoneType(createDNSRecordsetClient, d.Get("zone_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	// createDNSRecordset: create DNS recordset.
	if err := createDNSRecordset(createDNSRecordsetClient, d, zoneType); err != nil {
		return diag.FromErr(err)
	}

	zoneID, recordsetID, err := parseDNSRecordsetID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	waitForConfig := &WaitForConfig{
		ZoneID:      zoneID,
		RecordsetID: recordsetID,
		ZoneType:    zoneType,
		Timeout:     d.Timeout(schema.TimeoutCreate),
	}
	if err := waitForDNSRecordsetCreateOrUpdate(ctx, createDNSRecordsetClient, waitForConfig); err != nil {
		return diag.FromErr(err)
	}

	return resourceDNSRecordsetRead(ctx, d, meta)
}

func createDNSRecordset(recordsetClient *golangsdk.ServiceClient, d *schema.ResourceData, zoneType string) error {
	version := getApiVersionByZoneType(zoneType)
	createDNSRecordsetHttpUrl := fmt.Sprintf("%s/zones/{zone_id}/recordsets", version)

	zoneID := d.Get("zone_id").(string)
	createDNSRecordsetPath := recordsetClient.Endpoint + createDNSRecordsetHttpUrl
	createDNSRecordsetPath = strings.ReplaceAll(createDNSRecordsetPath, "{zone_id}", zoneID)

	createDNSRecordsetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	createDNSRecordsetOpt.JSONBody = utils.RemoveNil(buildCreateDNSRecordsetBodyParams(d))
	createDNSRecordsetResp, err := recordsetClient.Request("POST", createDNSRecordsetPath,
		&createDNSRecordsetOpt)
	if err != nil {
		return fmt.Errorf("error creating DNS recordset: %s", err)
	}

	createDNSRecordsetRespBody, err := utils.FlattenResponse(createDNSRecordsetResp)
	if err != nil {
		return err
	}

	id, err := jmespath.Search("id", createDNSRecordsetRespBody)
	if err != nil {
		return fmt.Errorf("error creating DNS recordset: ID is not found in API response")
	}
	d.SetId(fmt.Sprintf("%s/%s", zoneID, id))
	return nil
}

func waitForDNSRecordsetCreateOrUpdate(ctx context.Context, recordsetClient *golangsdk.ServiceClient,
	waitForConfig *WaitForConfig) error {
	stateConf := &resource.StateChangeConf{
		Target:       []string{"ACTIVE", "DISABLE"},
		Pending:      []string{"PENDING"},
		Refresh:      dnsRecordsetStatusRefreshFunc(recordsetClient, waitForConfig),
		Timeout:      waitForConfig.Timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for DNS recordset (%s) to be ACTIVE or DISABLE : %s",
			waitForConfig.RecordsetID, err)
	}
	return nil
}

func buildCreateDNSRecordsetBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"type":        utils.ValueIgnoreEmpty(d.Get("type")),
		"status":      utils.ValueIgnoreEmpty(d.Get("status")),
		"ttl":         utils.ValueIgnoreEmpty(d.Get("ttl")),
		"records":     utils.ValueIgnoreEmpty(d.Get("records")),
		"tags":        utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
	}
	return bodyParams
}

func resourceDNSRecordsetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDNSRecordset: Query DNS recordset
	getDNSRecordsetClient, err := cfg.NewServiceClient("dns_region", region)
	if err != nil {
		return diag.Errorf("error creating DNS Client: %s", err)
	}

	zoneID, recordsetID, err := parseDNSRecordsetID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	zoneType, err := getDNSZoneType(getDNSRecordsetClient, zoneID)
	if err != nil {
		return diag.FromErr(err)
	}
	version := getApiVersionByZoneType(zoneType)
	getDNSRecordsetHttpUrl := fmt.Sprintf("%s/zones/{zone_id}/recordsets/{recordset_id}", version)

	getDNSRecordsetPath := getDNSRecordsetClient.Endpoint + getDNSRecordsetHttpUrl
	getDNSRecordsetPath = strings.ReplaceAll(getDNSRecordsetPath, "{zone_id}", zoneID)
	getDNSRecordsetPath = strings.ReplaceAll(getDNSRecordsetPath, "{recordset_id}", recordsetID)

	getDNSRecordsetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getDNSRecordsetResp, err := getDNSRecordsetClient.Request("GET", getDNSRecordsetPath, &getDNSRecordsetOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DNS recordset")
	}

	getDNSRecordsetRespBody, err := utils.FlattenResponse(getDNSRecordsetResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getDNSRecordsetRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getDNSRecordsetRespBody, nil)),
		d.Set("zone_id", utils.PathSearch("zone_id", getDNSRecordsetRespBody, nil)),
		d.Set("zone_name", utils.PathSearch("zone_name", getDNSRecordsetRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getDNSRecordsetRespBody, nil)),
		d.Set("ttl", utils.PathSearch("ttl", getDNSRecordsetRespBody, nil)),
		d.Set("records", utils.PathSearch("records", getDNSRecordsetRespBody, nil)),
		d.Set("status", getDNSRecordsetStatus(getDNSRecordsetRespBody)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	// set tags
	if err := setDNSRecordsetTags(d, getDNSRecordsetClient, recordsetID, zoneType); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func setDNSRecordsetTags(d *schema.ResourceData, client *golangsdk.ServiceClient, id, zoneType string) error {
	resourceType, err := utils.GetDNSRecordSetTagType(zoneType)
	if err != nil {
		return err
	}
	return utils.SetResourceTagsToState(d, client, resourceType, id)
}

func getDNSRecordsetStatus(getDNSRecordsetRespBody interface{}) string {
	status := utils.PathSearch("status", getDNSRecordsetRespBody, "").(string)
	if status == "ACTIVE" {
		return "ENABLE"
	}
	return status
}

func resourceDNSRecordsetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("type") {
		return diag.Errorf("recordset action not permitted: Can not support change recordset type.")
	}
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)

	recordsetClient, err := cfg.NewServiceClient("dns_region", region)
	if err != nil {
		return diag.Errorf("error creating DNS Client: %s", err)
	}

	zoneID, recordsetID, err := parseDNSRecordsetID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	zoneType, err := getDNSZoneType(recordsetClient, zoneID)
	if err != nil {
		return diag.FromErr(err)
	}

	updateDNSRecordsetChanges := []string{
		"name",
		"description",
		"type",
		"ttl",
		"records",
	}
	if d.HasChanges(updateDNSRecordsetChanges...) {
		// updateDNSRecordset: Update DNS recordset
		if err := updateDNSRecordset(recordsetClient, d, zoneID, recordsetID, zoneType); err != nil {
			return diag.FromErr(err)
		}

		waitForConfig := &WaitForConfig{
			ZoneID:      zoneID,
			RecordsetID: recordsetID,
			ZoneType:    zoneType,
			Timeout:     d.Timeout(schema.TimeoutUpdate),
		}
		if err := waitForDNSRecordsetCreateOrUpdate(ctx, recordsetClient, waitForConfig); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("status") {
		// updateDNSRecordsetStatus: Update DNS recordset status
		if err := updateDNSRecordsetStatus(recordsetClient, d, recordsetID); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		resourceType, err := utils.GetDNSRecordSetTagType(zoneType)
		if err != nil {
			return diag.FromErr(err)
		}

		err = utils.UpdateResourceTags(recordsetClient, d, resourceType, recordsetID)
		if err != nil {
			return diag.Errorf("error updating DNS recordset tags: %s", err)
		}
	}
	return resourceDNSRecordsetRead(ctx, d, meta)
}

func updateDNSRecordset(recordsetClient *golangsdk.ServiceClient, d *schema.ResourceData, zoneID,
	recordsetID, zoneType string) error {
	version := getApiVersionByZoneType(zoneType)
	updateDNSRecordsetHttpUrl := fmt.Sprintf("%s/zones/{zone_id}/recordsets/{recordset_id}", version)

	updateDNSRecordsetPath := recordsetClient.Endpoint + updateDNSRecordsetHttpUrl
	updateDNSRecordsetPath = strings.ReplaceAll(updateDNSRecordsetPath, "{zone_id}", zoneID)
	updateDNSRecordsetPath = strings.ReplaceAll(updateDNSRecordsetPath, "{recordset_id}", recordsetID)

	updateDNSRecordsetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	updateDNSRecordsetOpt.JSONBody = utils.RemoveNil(buildUpdateDNSRecordsetBodyParams(d))
	_, err := recordsetClient.Request("PUT", updateDNSRecordsetPath, &updateDNSRecordsetOpt)
	if err != nil {
		return fmt.Errorf("error updating DNS recordset: %s", err)
	}
	return nil
}

func updateDNSRecordsetStatus(recordsetClient *golangsdk.ServiceClient, d *schema.ResourceData,
	recordsetID string) error {
	var (
		updateDNSRecordsetStatusHttpUrl = "v2.1/recordsets/{recordset_id}/statuses/set"
	)

	updateDNSRecordsetStatusPath := recordsetClient.Endpoint + updateDNSRecordsetStatusHttpUrl
	updateDNSRecordsetStatusPath = strings.ReplaceAll(updateDNSRecordsetStatusPath, "{recordset_id}", recordsetID)

	updateDNSRecordsetStatusOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	updateDNSRecordsetStatusOpt.JSONBody = utils.RemoveNil(buildUpdateDNSRecordsetStatusBodyParams(d))
	_, err := recordsetClient.Request("PUT", updateDNSRecordsetStatusPath, &updateDNSRecordsetStatusOpt)
	if err != nil {
		return fmt.Errorf("error updating DNS recordset status: %s", err)
	}
	return nil
}

func buildUpdateDNSRecordsetBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"type":        utils.ValueIgnoreEmpty(d.Get("type")),
		"ttl":         utils.ValueIgnoreEmpty(d.Get("ttl")),
		"records":     utils.ValueIgnoreEmpty(d.Get("records")),
	}
	return bodyParams
}

func buildUpdateDNSRecordsetStatusBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"status": utils.ValueIgnoreEmpty(d.Get("status")),
	}
	return bodyParams
}

func resourceDNSRecordsetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)

	// deleteDNSRecordset: Delete DNS recordset
	deleteDNSRecordsetClient, err := cfg.NewServiceClient("dns_region", region)
	if err != nil {
		return diag.Errorf("error creating DNS Client: %s", err)
	}

	zoneID, recordsetID, err := parseDNSRecordsetID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	zoneType, err := getDNSZoneType(deleteDNSRecordsetClient, zoneID)
	if err != nil {
		return diag.FromErr(err)
	}
	version := getApiVersionByZoneType(zoneType)
	deleteDNSRecordsetHttpUrl := fmt.Sprintf("%s/zones/{zone_id}/recordsets/{recordset_id}", version)

	deleteDNSRecordsetPath := deleteDNSRecordsetClient.Endpoint + deleteDNSRecordsetHttpUrl
	deleteDNSRecordsetPath = strings.ReplaceAll(deleteDNSRecordsetPath, "{zone_id}", zoneID)
	deleteDNSRecordsetPath = strings.ReplaceAll(deleteDNSRecordsetPath, "{recordset_id}", recordsetID)

	deleteDNSRecordsetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	_, err = deleteDNSRecordsetClient.Request("DELETE", deleteDNSRecordsetPath, &deleteDNSRecordsetOpt)
	if err != nil {
		return diag.Errorf("error deleting DNS recordset: %s", err)
	}

	waitForConfig := &WaitForConfig{
		ZoneID:      zoneID,
		RecordsetID: recordsetID,
		ZoneType:    zoneType,
		Timeout:     d.Timeout(schema.TimeoutDelete),
	}
	if err := waitForDNSRecordsetDeleted(ctx, deleteDNSRecordsetClient, waitForConfig); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func waitForDNSRecordsetDeleted(ctx context.Context, recordsetClient *golangsdk.ServiceClient,
	waitForConfig *WaitForConfig) error {
	stateConf := &resource.StateChangeConf{
		Target:       []string{"DELETED"},
		Pending:      []string{"ACTIVE", "PENDING", "ERROR"},
		Refresh:      dnsRecordsetStatusRefreshFunc(recordsetClient, waitForConfig),
		Timeout:      waitForConfig.Timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for DNS recordset (%s) to be DELETED: %s",
			waitForConfig.RecordsetID, err)
	}
	return nil
}

func dnsRecordsetStatusRefreshFunc(client *golangsdk.ServiceClient, waitForConfig *WaitForConfig) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		version := getApiVersionByZoneType(waitForConfig.ZoneType)
		getDNSRecordsetHttpUrl := fmt.Sprintf("%s/zones/{zone_id}/recordsets/{recordset_id}", version)

		getDNSRecordsetPath := client.Endpoint + getDNSRecordsetHttpUrl
		getDNSRecordsetPath = strings.ReplaceAll(getDNSRecordsetPath, "{zone_id}", waitForConfig.ZoneID)
		getDNSRecordsetPath = strings.ReplaceAll(getDNSRecordsetPath, "{recordset_id}", waitForConfig.RecordsetID)

		getDNSRecordsetOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		getDNSRecordsetResp, err := client.Request("GET", getDNSRecordsetPath, &getDNSRecordsetOpt)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return getDNSRecordsetResp, "DELETED", nil
			}
			return nil, "", err
		}

		getDNSRecordsetRespBody, err := utils.FlattenResponse(getDNSRecordsetResp)
		if err != nil {
			return nil, "", err
		}

		status := utils.PathSearch("status", getDNSRecordsetRespBody, "")
		return getDNSRecordsetRespBody, parseStatus(status.(string)), nil
	}
}

func parseStatus(rawStatus string) string {
	splits := strings.Split(rawStatus, "_")
	// rawStatus maybe one of PENDING_CREATE, PENDING_UPDATE, PENDING_DELETE, ACTIVE, or ERROR
	return splits[0]
}

func parseDNSRecordsetID(id string) (zoneID, recordsetID string, err error) {
	idArrays := strings.SplitN(id, "/", 2)
	if len(idArrays) != 2 {
		err = fmt.Errorf("invalid format specified for ID. Format must be <zone_id>/<recordset_id>")
		return
	}
	zoneID = idArrays[0]
	recordsetID = idArrays[1]
	return
}

func getApiVersionByZoneType(zoneType string) string {
	if zoneType == "private" {
		return "v2"
	}
	return "v2.1"
}

func getDNSZoneType(dnsClient *golangsdk.ServiceClient, zoneID string) (string, error) {
	zoneInfo, err := zones.Get(dnsClient, zoneID).Extract()
	if err != nil {
		return "", fmt.Errorf("error getting zone: %s", err)
	}
	return zoneInfo.ZoneType, nil
}
