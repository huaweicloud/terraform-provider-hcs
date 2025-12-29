package lts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

// LTS GET /v2/{project_id}/log-streams
func DataSourceLogStreams() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLogStreamsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the log streams.`,
			},
			"log_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the log group to be queried.`,
			},
			"log_stream_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the log stream to be queried.`,
			},
			"log_streams": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of log streams.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the log stream.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the log stream.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the log stream, in RFC3339 format.`,
						},
						"filter_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of filters.`,
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: `The tags of the log stream.`,
						},
						"is_favorite": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the stream is bookmarked.`,
						},
					},
				},
			},
		},
	}
}

func buildLogStreamsQueryParams(d *schema.ResourceData) string {
	res := ""
	if logGroupName, ok := d.GetOk("log_group_name"); ok {
		res += "&log_group_name=" + logGroupName.(string)
	}
	if logStreamName, ok := d.GetOk("log_stream_name"); ok {
		res += "&log_stream_name=" + logStreamName.(string)
	}

	return res
}

func getLogStreams(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/log-streams?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{limit}", strconv.Itoa(limit))
	getPath += buildLogStreamsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	for {
		getPathWithOffset := getPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", getPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		logStreams := utils.PathSearch("log_streams", respBody, make([]interface{}, 0)).([]interface{})
		if logStreams == nil {
			break
		}

		result = append(result, logStreams...)

		if len(logStreams) < limit {
			break
		}
		offset += len(logStreams)
	}

	return result, nil
}

func flattenLogStreamsList(logStreams []interface{}) []map[string]interface{} {
	if len(logStreams) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(logStreams))
	for _, logStream := range logStreams {
		creationTime := utils.PathSearch("creation_time", logStream, float64(0)).(float64) / 1000
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("log_stream_id", logStream, nil),
			"name":         utils.PathSearch("log_stream_name", logStream, nil),
			"created_at":   utils.FormatTimeStampRFC3339(int64(creationTime), false),
			"filter_count": utils.PathSearch("filter_count", logStream, nil),
			"is_favorite":  utils.PathSearch("is_favorite", logStream, nil),
			"tags":         utils.DeleteEnterpriseProjectIdFromTags(utils.PathSearch("tag", logStream, nil)),
		})
	}

	return result
}

func dataSourceLogStreamsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	logStreams, err := getLogStreams(client, d)
	if err != nil {
		return diag.Errorf("error getting log streams: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("log_streams", flattenLogStreamsList(logStreams)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
