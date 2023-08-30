package smn

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/common/tags"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/smn/v2/topics"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceTopic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTopicCreate,
		ReadContext:   resourceTopicRead,
		UpdateContext: resourceTopicUpdate,
		DeleteContext: resourceTopicDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9][\w-]*$`),
						"The name must start with a letter or digit, and can only contain "+
							"letters, digits, underscores (_), and hyphens (-)."),
					validation.StringLenBetween(1, 255),
				),
			},
			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 192),
			},
			"tags": common.TagsSchema(),

			"topic_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"push_policy": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTopicCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	createOpts := topics.CreateOps{
		Name:                d.Get("name").(string),
		DisplayName:         d.Get("display_name").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}

	log.Printf("[DEBUG] create Options: %#v", createOpts)
	topic, err := topics.Create(client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error getting SMN topic from result: %s", err)
	}

	log.Printf("[DEBUG] successfully created SMN topic: %s", topic.TopicUrn)
	d.SetId(topic.TopicUrn)

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		tagClient, err := cfg.SmnV2TagClient(region)
		if err != nil {
			return diag.Errorf("error creating SMN tag client: %s", err)
		}
		tagClient.MoreHeaders = map[string]string{
			"X-SMN-RESOURCEID-TYPE": "name",
		}
		if tagErr := tags.Create(tagClient, "smn_topic", d.Get("name").(string), taglist).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of SMN topic %s: %s", d.Id(), tagErr)
		}
	}

	return resourceTopicRead(ctx, d, meta)
}

func resourceTopicRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	topicUrn := d.Id()
	topicGet, err := topics.Get(client, topicUrn).ExtractGet()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SMN topic")
	}

	log.Printf("[DEBUG] retrieved SMN topic %s: %#v", topicUrn, topicGet)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("topic_urn", topicGet.TopicUrn),
		d.Set("display_name", topicGet.DisplayName),
		d.Set("name", topicGet.Name),
		d.Set("push_policy", topicGet.PushPolicy),
		d.Set("update_time", topicGet.UpdateTime),
		d.Set("create_time", topicGet.CreateTime),
	)

	// fetch tags
	tagClient, err := cfg.SmnV2TagClient(region)
	if err != nil {
		return diag.Errorf("error creating SMN tag client: %s", err)
	}
	tagClient.MoreHeaders = map[string]string{
		"X-SMN-RESOURCEID-TYPE": "name",
	}
	if resourceTags, err := tags.Get(tagClient, "smn_topic", d.Get("name").(string)).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		mErr = multierror.Append(mErr, d.Set("tags", tagmap))
	} else {
		log.Printf("[WARN] fetching tags of SMN topic failed: %s", err)
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting SMN topic fields: %s", err)
	}

	return nil
}

func resourceTopicUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	id := d.Id()
	if d.HasChange("display_name") {
		updateOpts := topics.UpdateOps{
			DisplayName: d.Get("display_name").(string),
		}
		if _, err = topics.Update(client, updateOpts, id).Extract(); err != nil {
			return diag.Errorf("error updating SMN topic %s: %s", id, err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		tagClient, err := cfg.SmnV2TagClient(region)
		if err != nil {
			return diag.Errorf("error creating SMN tag client: %s", err)
		}
		tagClient.MoreHeaders = map[string]string{
			"X-SMN-RESOURCEID-TYPE": "name",
		}
		tagErr := utils.UpdateResourceTags(tagClient, d, "smn_topic", d.Get("name").(string))
		if tagErr != nil {
			return diag.Errorf("error updating tags of SMN topic %s: %s", id, tagErr)
		}
	}

	return resourceTopicRead(ctx, d, meta)
}

func resourceTopicDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	client, err := cfg.SmnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	result := topics.Delete(client, d.Id())
	if result.Err != nil {
		return diag.Errorf("error deleting SMN topic: %s", result.Err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForTopicDelete(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting SMN topic %s: %s", d.Id(), err)
	}

	return nil
}

func waitForTopicDelete(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := topics.Get(client, id).ExtractGet()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] successfully deleted topic %s", id)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}
		return r, "ACTIVE", nil
	}
}
