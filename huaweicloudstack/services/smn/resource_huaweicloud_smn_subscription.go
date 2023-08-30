package smn

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/smn/v2/subscriptions"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubscriptionCreate,
		ReadContext:   resourceSubscriptionRead,
		DeleteContext: resourceSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSubscriptionImport,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"topic_urn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"email", "sms",
				}, false),
			},
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"subscription_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceSubscriptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	topicUrn := d.Get("topic_urn").(string)
	createOpts := subscriptions.CreateOps{
		Endpoint: d.Get("endpoint").(string),
		Protocol: d.Get("protocol").(string),
		Remark:   d.Get("remark").(string),
	}

	log.Printf("[DEBUG] create Options: %#v", createOpts)
	subscription, err := subscriptions.Create(client, createOpts, topicUrn).Extract()
	if err != nil {
		return diag.Errorf("error getting subscription from result: %s", err)
	}

	log.Printf("[DEBUG] create SMN subscription: %s", subscription.SubscriptionUrn)
	d.SetId(subscription.SubscriptionUrn)
	return resourceSubscriptionRead(ctx, d, meta)
}

func resourceSubscriptionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	id := d.Id()
	topicUrn := d.Get("topic_urn").(string)

	log.Printf("[DEBUG] fetching subscription: %s", id)
	subscriptionslist, err := subscriptions.ListFromTopic(client, topicUrn).Extract()
	if err != nil {
		return diag.Errorf("error fetching the list of subscriptions: %s", err)
	}

	var targetSubscription *subscriptions.SubscriptionGet
	for i := range subscriptionslist {
		if subscriptionslist[i].SubscriptionUrn == id {
			targetSubscription = &subscriptionslist[i]
			break
		}
	}

	if targetSubscription == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("topic_urn", targetSubscription.TopicUrn),
		d.Set("endpoint", targetSubscription.Endpoint),
		d.Set("protocol", targetSubscription.Protocol),
		d.Set("subscription_urn", targetSubscription.SubscriptionUrn),
		d.Set("owner", targetSubscription.Owner),
		d.Set("remark", targetSubscription.Remark),
		d.Set("status", targetSubscription.Status),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting SMN topic fields: %s", err)
	}
	return nil
}

func resourceSubscriptionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	id := d.Id()
	topicUrn := d.Get("topic_urn").(string)
	log.Printf("[DEBUG] successfully delete topicUrn %s", topicUrn)
	result := subscriptions.Delete(client, url.QueryEscape(topicUrn), url.QueryEscape(id))
	if result.Err != nil {
		return diag.FromErr(result.Err)
	}

	log.Printf("[DEBUG] successfully delete subscription %s", id)
	return nil
}

func resourceSubscriptionImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	subscriptionUrn := d.Id()
	index := strings.LastIndex(subscriptionUrn, ":")
	if index == -1 {
		return nil, fmt.Errorf("invalid format: the subscription URN is invalid")
	}
	topicUrn := subscriptionUrn[:index]

	d.SetId(subscriptionUrn)
	d.Set("topic_urn", topicUrn)

	return []*schema.ResourceData{d}, nil
}
