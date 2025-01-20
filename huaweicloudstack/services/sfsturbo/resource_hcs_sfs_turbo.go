package sfsturbo

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/sfs_turbo/v1/shares"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

const (
	prepaidUnitMonth int = 2
	prepaidUnitYear  int = 3

	autoRenewDisabled int = 0
	autoRenewEnabled  int = 1

	shareTypeSsd = "sfsturbo.ssd"
	shareTypeHdd = "sfsturbo.hdd"
)

func ResourceSFSTurbo() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSFSTurboCreate,
		ReadContext:   resourceSFSTurboRead,
		UpdateContext: resourceSFSTurboUpdate,
		DeleteContext: resourceSFSTurboDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"share_proto": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "NFS",
				ValidateFunc: validation.StringInSlice([]string{"NFS"}, false),
			},
			"share_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  shareTypeHdd,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			},
			"dedicated_flavor": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dedicated_storage_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"export_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_capacity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),
		},
	}
}

func buildTurboMetadataOpts(d *schema.ResourceData) shares.Metadata {
	metaOpts := shares.Metadata{}
	if v, ok := d.GetOk("dedicated_flavor"); ok {
		metaOpts.DedicatedFlavor = v.(string)
	}
	if v, ok := d.GetOk("dedicated_storage_id"); ok {
		metaOpts.DedicatedStorageID = v.(string)
	}
	return metaOpts
}

func buildTurboCreateOpts(cfg *config.HcsConfig, d *schema.ResourceData) shares.CreateOpts {
	result := shares.CreateOpts{
		Share: shares.Share{
			Name:                d.Get("name").(string),
			Size:                d.Get("size").(int),
			Bandwidth:           d.Get("bandwidth").(int),
			ShareProto:          d.Get("share_proto").(string),
			VpcID:               d.Get("vpc_id").(string),
			SubnetID:            d.Get("subnet_id").(string),
			SecurityGroupID:     d.Get("security_group_id").(string),
			AvailabilityZone:    d.Get("availability_zone").(string),
			ShareType:           d.Get("share_type").(string),
			EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
			Metadata:            buildTurboMetadataOpts(d),
		},
	}

	if d.Get("charging_mode") == "prePaid" {
		billing := shares.BssParam{
			PeriodNum: d.Get("period").(int),
			IsAutoPay: utils.Int(1), // Always enable auto-pay.
		}
		if d.Get("period_unit").(string) == "month" {
			billing.PeriodType = prepaidUnitMonth
		} else {
			billing.PeriodType = prepaidUnitYear
		}
		if d.Get("auto_renew").(string) == "true" {
			billing.IsAutoRenew = utils.Int(autoRenewEnabled)
		} else {
			billing.IsAutoRenew = utils.Int(autoRenewDisabled)
		}
		result.BssParam = &billing
	}

	return result
}

func resourceSFSTurboCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	sfsClient, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	createOpts := buildTurboCreateOpts(cfg, d)
	log.Printf("[DEBUG] create sfs turbo with option: %+v", createOpts)
	resp, err := shares.Create(sfsClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating SFS Turbo: %s", err)
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		orderId := resp.OrderId
		if orderId == "" {
			return diag.Errorf("unable to find the order ID, this is a COM (Cloud Order Management) error, " +
				"please contact service for help and check your order status on the console.")
		}
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(resourceId)
	} else {
		d.SetId(resp.ID)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"100"},
		Target:       []string{"200"},
		Refresh:      waitForSFSTurboStatus(sfsClient, resp.ID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		PollInterval: 3 * time.Second,
	}
	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return diag.Errorf("error waiting for SFS Turbo (%s) to become ready: %s ", d.Id(), stateErr)
	}

	// add tags
	if err := utils.CreateResourceTags(sfsClient, d, "sfs-turbo", d.Id()); err != nil {
		return diag.Errorf("error setting tags of SFS Turbo %s: %s", d.Id(), err)
	}

	return resourceSFSTurboRead(ctx, d, meta)
}

func flattenSize(n *shares.Turbo) interface{} {
	// n.Size is a string of float64, should convert it to int
	if fsize, err := strconv.ParseFloat(n.Size, 64); err == nil {
		return int(fsize)
	}

	return nil
}

func flattenStatus(n *shares.Turbo) interface{} {
	if n.SubStatus != "" {
		return n.SubStatus
	}

	return n.Status
}

func resourceSFSTurboRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	sfsClient, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	n, err := shares.Get(sfsClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "SFS Turbo")
	}

	mErr := multierror.Append(
		nil,
		d.Set("name", n.Name),
		d.Set("share_proto", n.ShareProto),
		d.Set("share_type", n.ShareType),
		d.Set("vpc_id", n.VpcID),
		d.Set("subnet_id", n.SubnetID),
		d.Set("security_group_id", n.SecurityGroupID),
		d.Set("version", n.Version),
		d.Set("region", region),
		d.Set("availability_zone", n.AvailabilityZone),
		d.Set("available_capacity", n.AvailCapacity),
		d.Set("export_location", n.ExportLocation),
		d.Set("enterprise_project_id", n.EnterpriseProjectId),
		d.Set("bandwidth", n.Bandwidth),
		d.Set("size", flattenSize(n)),
		d.Set("status", flattenStatus(n)),
	)

	// set tags
	err = utils.SetResourceTagsToState(d, sfsClient, "sfs-turbo", d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diag.FromErr(mErr.ErrorOrNil())
}

// buildTurboUpdateOpts supports SFS-Turbo size and bandwidth(bandwidth is HCS specific)
func buildTurboUpdateOpts(oldSize, newSize, oldBandwidth, newBandwidth int) shares.ExpandOpts {
	expandOpts := shares.ExtendOpts{}

	expandOpts.BssParam = &shares.BssParamExtend{
		IsAutoPay: utils.Int(1),
	}

	if oldSize == newSize && oldBandwidth != newBandwidth {
		expandOpts = shares.ExtendOpts{
			NewBandwidth: newBandwidth,
		}
	}

	if oldSize != newSize && oldBandwidth == newBandwidth {
		expandOpts = shares.ExtendOpts{
			NewSize: newSize,
		}
	}

	// This circumstance the API will return error: can not be changed at the same time.
	if oldSize != newSize && oldBandwidth != newBandwidth {
		expandOpts = shares.ExtendOpts{
			NewSize:      newSize,
			NewBandwidth: newBandwidth,
		}
	}

	return shares.ExpandOpts{
		Extend: expandOpts,
	}
}

func resourceSFSTurboUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	sfsClient, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	resourceId := d.Id()
	if d.HasChanges("size", "bandwidth") {
		// size and bandwidth
		oldBandwidth, newBandwidth := d.GetChange("bandwidth")
		oldSize, newSize := d.GetChange("size")
		if oldSize.(int) > newSize.(int) {
			return diag.Errorf("shrinking SFS Turbo size is not supported")
		}

		updateOpts := buildTurboUpdateOpts(oldSize.(int), newSize.(int), oldBandwidth.(int), newBandwidth.(int))
		resp, err := shares.Expand(sfsClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error expanding SFS Turbo size: %s", err)
		}

		// charging_mode
		isPrePaid := d.Get("charging_mode").(string) == "prePaid"
		if isPrePaid {
			orderId := resp.OrderId
			if orderId == "" {
				return diag.Errorf("unable to find the order ID, this is a COM (Cloud Order Management) error, " +
					"please contact service for help and check your order status on the console.")
			}
			bssClient, err := cfg.BssV2Client(region)
			if err != nil {
				return diag.Errorf("error creating BSS v2 client: %s", err)
			}
			err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
			_, err = common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		stateConf := &resource.StateChangeConf{
			Pending:      []string{"121"},
			Target:       []string{"221", "200"},
			Refresh:      waitForSFSTurboSubStatus(sfsClient, resourceId),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			PollInterval: 5 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error updating SFS Turbo: %s", err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		if err := updateSFSTurboTags(sfsClient, d); err != nil {
			return diag.Errorf("error updating tags of SFS Turbo %s: %s", resourceId, err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), resourceId); err != nil {
			return diag.Errorf("error updating the auto-renew of the SFS Turbo (%s): %s", resourceId, err)
		}
	}

	if d.HasChange("name") {
		updateNameOpts := shares.UpdateNameOpts{
			Name: d.Get("name").(string),
		}
		err = shares.UpdateName(sfsClient, d.Id(), updateNameOpts).Err
		if err != nil {
			return diag.Errorf("error updating name of SFS Turbo: %s", err)
		}
	}

	if d.HasChange("security_group_id") {
		updateSecurityGroupIdOpts := shares.UpdateSecurityGroupIdOpts{
			SecurityGroupId: d.Get("security_group_id").(string),
		}
		err = shares.UpdateSecurityGroupId(sfsClient, d.Id(), updateSecurityGroupIdOpts).Err
		if err != nil {
			return diag.Errorf("error updating security group ID of SFS Turbo: %s", err)
		}
		stateConf := &resource.StateChangeConf{
			Pending:      []string{"132"},
			Target:       []string{"232", "200"},
			Refresh:      waitForSFSTurboSubStatus(sfsClient, resourceId),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			PollInterval: 5 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error updating SFS Turbo: %s", err)
		}
	}

	return resourceSFSTurboRead(ctx, d, meta)
}

func updateSFSTurboTags(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	// remove old tags
	oldKeys := getOldTagKeys(d)
	if err := utils.DeleteResourceTagsWithKeys(client, oldKeys, "sfs-turbo", d.Id()); err != nil {
		return err
	}

	// set new tags
	return utils.CreateResourceTags(client, d, "sfs-turbo", d.Id())
}

func getOldTagKeys(d *schema.ResourceData) []string {
	oRaw, _ := d.GetChange("tags")
	var tagKeys []string
	if oMap := oRaw.(map[string]interface{}); len(oMap) > 0 {
		for k := range oMap {
			tagKeys = append(tagKeys, k)
		}
	}
	return tagKeys
}

func resourceSFSTurboDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	sfsClient, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	resourceId := d.Id()
	// for prePaid mode, we should unsubscribe the resource
	if d.Get("charging_mode").(string) == "prePaid" {
		err := common.UnsubscribePrePaidResource(d, cfg, []string{resourceId})
		if err != nil {
			return diag.Errorf("error unsubscribing SFS Turbo: %s", err)
		}
	} else {
		err = shares.Delete(sfsClient, resourceId).ExtractErr()
		if err != nil {
			return common.CheckDeletedDiag(d, err, "SFS Turbo")
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"100", "200"},
		Target:     []string{"deleted"},
		Refresh:    waitForSFSTurboStatus(sfsClient, resourceId),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting SFS Turbo: %s", err)
	}
	return nil
}

func waitForSFSTurboStatus(sfsClient *golangsdk.ServiceClient, shareId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := shares.Get(sfsClient, shareId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted shared File %s", shareId)
				return r, "deleted", nil
			}
			return r, "error", err
		}

		return r, r.Status, nil
	}
}

func waitForSFSTurboSubStatus(sfsClient *golangsdk.ServiceClient, shareId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := shares.Get(sfsClient, shareId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted shared File %s", shareId)
				return r, "deleted", nil
			}
			return r, "error", err
		}

		var status string
		if r.SubStatus != "" {
			status = r.SubStatus
		} else {
			status = r.Status
		}
		return r, status, nil
	}
}
