package lts

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

// LTS GET /v2/{project_id}/groups
func DataSourceLogGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLogGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the log groups.`,
			},
			"log_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of log groups.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the log group.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the log group.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the log group, in RFC3339 format.`,
						},
						"ttl_in_days": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The storage duration of the log group in days.`,
						},
						"stream_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of log streams under the log group.`,
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: `The tags of the log group.`,
						},
					},
				},
			},
		},
	}
}

func getLogGroups(client *golangsdk.ServiceClient, _ *schema.ResourceData) ([]interface{}, error) {
	var httpUrl = "v2/{project_id}/groups"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
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

	logGroups := utils.PathSearch("log_groups", respBody, make([]interface{}, 0)).([]interface{})
	if len(logGroups) == 0 {
		return make([]interface{}, 0), nil
	}

	return logGroups, nil
}

func flattenLogGroups(logGroups []interface{}) []map[string]interface{} {
	if len(logGroups) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(logGroups))
	for _, logGroup := range logGroups {
		creationTime := utils.PathSearch("creation_time", logGroup, float64(0)).(float64) / 1000
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("log_group_id", logGroup, nil),
			"name":        utils.PathSearch("log_group_name", logGroup, nil),
			"ttl_in_days": utils.PathSearch("ttl_in_days", logGroup, nil),
			"stream_size": utils.PathSearch("stream_size", logGroup, nil),
			"tags":        utils.DeleteEnterpriseProjectIdFromTags(utils.PathSearch("tag", logGroup, nil)),
			"created_at":  utils.FormatTimeStampRFC3339(int64(creationTime), false),
		})
	}

	return result
}

func dataSourceLogGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	logGroups, err := getLogGroups(client, d)
	if err != nil {
		return diag.Errorf("error getting log groups: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("log_groups", flattenLogGroups(logGroups)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
