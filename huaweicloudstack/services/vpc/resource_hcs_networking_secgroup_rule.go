package vpc

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	v1Rules "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/security/rules"
	v3Rules "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v3/security/rules"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/logp"
)

// Some parameters are only support creation in ver.3 API.
var advancedParams = []string{"ports", "action"}

func ResourceNetworkingSecGroupRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingSecGroupRuleCreate,
		ReadContext:   resourceNetworkingSecGroupRuleRead,
		DeleteContext: resourceNetworkingSecGroupRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ingress", "egress",
				}, true),
			},
			"ethertype": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"IPv4", "IPv6",
				}, true),
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			// The port range parameters conflict with advanced parameters.
			"port_range_min": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				RequiredWith: []string{"protocol"},
			},
			"port_range_max": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				RequiredWith: []string{"port_range_min"},
			},
			"ports": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Computed:   true,
				Deprecated: "use port_range_min and port_range_max instead",
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.Any(
					validation.StringInSlice([]string{"tcp", "udp", "icmp", "icmpv6"}, false),
					validation.StringMatch(regexp.MustCompile("^([0-1]?[0-9]?[0-9]|2[0-4][0-9]|25[0-5])$"),
						"The valid protocol is range from 0 to 255.",
					),
				),
			},
			"remote_group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ExactlyOneOf: []string{"remote_ip_prefix"},
			},
			"remote_ip_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: utils.ValidateCIDR,
				StateFunc: func(v interface{}) string {
					return strings.ToLower(v.(string))
				},
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"allow", "deny",
				}, false),
				ConflictsWith: []string{"port_range_min", "port_range_max"},
			},
			"priority": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.IntBetween(1, 100),
				ConflictsWith: []string{"port_range_min", "port_range_max"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkingSecGroupRuleCreateV1(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	v1Client, err := cfg.NetworkingV1Client(common.GetRegion(d, cfg))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloudStack networking v1 client: %s", err)
	}

	opt := v1Rules.CreateOpts{
		Description:     d.Get("description").(string),
		SecurityGroupId: d.Get("security_group_id").(string),
		RemoteGroupId:   d.Get("remote_group_id").(string),
		RemoteIpPrefix:  d.Get("remote_ip_prefix").(string),
		Protocol:        d.Get("protocol").(string),
		Ethertype:       d.Get("ethertype").(string),
		Direction:       d.Get("direction").(string),
		PortRangeMin:    d.Get("port_range_min").(int),
		PortRangeMax:    d.Get("port_range_max").(int),
	}

	logp.Printf("[DEBUG] The createOpts of the Security Group rule is: %#v", opt)
	resp, err := v1Rules.Create(v1Client, opt)
	if err != nil {
		return fmtp.DiagErrorf("Error creating Security Group rule: %s", err)
	}
	d.SetId(resp.ID)

	return resourceNetworkingSecGroupRuleRead(ctx, d, meta)
}

func resourceNetworkingSecGroupRuleCreate(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	return resourceNetworkingSecGroupRuleCreateV1(ctx, d, meta)
}

func resourceNetworkingSecGroupRuleRead(_ context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := common.GetRegion(d, cfg)

	v1Client, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloudStack networking v1 client: %s", err)
	}
	v3Client, err := cfg.NetworkingV3Client(common.GetRegion(d, cfg))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloudStack networking v3 client: %s", err)
	}

	resp, err := v1Rules.Get(v1Client, d.Id())
	if err != nil {
		logp.Printf("[DEBUG] Unable to find the specified Security group rule (%s).", d.Id())
		return common.CheckDeletedDiag(d, err, "HuaweiCloudStack Security Group Rule")
	}

	mErr := multierror.Append(nil,
		d.Set("region", common.GetRegion(d, cfg)),
		d.Set("direction", resp.Direction),
		d.Set("description", resp.Description),
		d.Set("ethertype", resp.Ethertype),
		d.Set("protocol", resp.Protocol),
		d.Set("remote_group_id", resp.RemoteGroupId),
		d.Set("remote_ip_prefix", resp.RemoteIpPrefix),
		d.Set("security_group_id", resp.SecurityGroupId),
		d.Set("port_range_min", resp.PortRangeMin),
		d.Set("port_range_max", resp.PortRangeMax),
	)

	rule, err := v3Rules.Get(v3Client, d.Id())
	if err == nil {
		// If the v3 API method has no error, parse its ports attribute and setup.
		logp.Printf("[DEBUG] Retrieved Security Group rule (%s): %+v", d.Id(), rule)
		mErr = multierror.Append(mErr,
			d.Set("ports", rule.MultiPort),
			d.Set("action", rule.Action),
			d.Set("priority", rule.Priority),
		)
	}

	// If the query process returns an error, either because the specified region does not exist or the v3 API is
	// not released, or other reasons, skip the setting.
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceNetworkingSecGroupRuleDelete(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	logp.Printf("[DEBUG] Destroy security group rule: %s", d.Id())

	cfg := config.GetHcsConfig(meta)
	client, err := cfg.NetworkingV1Client(common.GetRegion(d, cfg))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloudStack networking v1 client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForSecGroupRuleDelete(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      8 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf("Error deleting HuaweiCloudStack Security Group Rule: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForSecGroupRuleDelete(client *golangsdk.ServiceClient, ruleId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		logp.Printf("[DEBUG] Attempting to delete HuaweiCloudStack Security Group Rule %s.", ruleId)

		r, err := v1Rules.Get(client, ruleId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloudStack Security Group Rule %s", ruleId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = v1Rules.Delete(client, ruleId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloudStack Security Group Rule %s", ruleId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		logp.Printf("[DEBUG] HuaweiCloudStack Security Group Rule %s still active.", ruleId)
		return r, "ACTIVE", nil
	}
}
