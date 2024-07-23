package vpc

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/flowlogs"
)

func DataSourceVpcFlowLog() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcFlowLogsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the VPC flow log name.`,
			},
			"flow_log_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the VPC flow log ID.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the resource type for which that the logs to be collected.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the resource ID for which that the logs to be collected.`,
			},
			"log_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the LTS log group ID.`,
			},
			"log_stream_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the LTS log stream ID.`,
			},
			"traffic_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of traffic to log.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status of the flow log.`,
			},
			"flow_logs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of VPC flow logs.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The VPC flow log name.`,
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of a VPC flow log`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The VPC flow log description.`,
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource type for which that the logs to be collected.`,
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource ID for which that the logs to be collected.`,
						},
						"log_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The LTS log group ID.`,
						},
						"log_stream_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The LTS log stream ID.`,
						},
						"traffic_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of traffic to log.`,
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether to enable the VPC flow log.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The VPC flow log status.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the resource is created.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the resource is last updated.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceVpcFlowLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	region := hcsConfig.GetRegion(d)
	client, err := hcsConfig.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	listOpts := flowlogs.ListOpts{
		ID:           d.Get("flow_log_id").(string),
		Name:         d.Get("name").(string),
		ResourceType: d.Get("resource_type").(string),
		ResourceID:   d.Get("resource_id").(string),
		TrafficType:  d.Get("traffic_type").(string),
		LogGroupID:   d.Get("log_group_id").(string),
		LogTopicID:   d.Get("log_stream_id").(string),
		Status:       d.Get("status").(string),
	}

	pages, err := flowlogs.List(client, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("unable to retrieve vpc flowlogs: %s", err)
	}

	body, err := flowlogs.ExtractFlowLogs(pages)
	if err != nil {
		return diag.Errorf("extract vpc flowlogs: %s", err)
	}

	log.Printf("[DEBUG] Retrieved flowlogs using given filter: %+v", body)

	var fls []map[string]interface{}
	var ids []string
	for _, item := range body {
		fl := map[string]interface{}{
			"name":          item.Name,
			"id":            item.ID,
			"description":   item.Description,
			"resource_type": item.ResourceType,
			"resource_id":   item.ResourceID,
			"log_group_id":  item.LogGroupID,
			"log_stream_id": item.LogTopicID,
			"traffic_type":  item.TrafficType,
			"enabled":       item.AdminState,
			"status":        item.Status,
			"created_at":    item.CreatedAt,
			"updated_at":    item.UpdatedAt,
		}

		fls = append(fls, fl)
		ids = append(ids, item.ID)
	}
	log.Printf("[DEBUG] vpc flowlogs List after filter, count=%d :%+v", len(fls), fls)

	mErr := d.Set("flow_logs", fls)
	if mErr != nil {
		return diag.Errorf("set flowlogs err:%s", mErr)
	}

	d.SetId(hashcode.Strings(ids))
	return nil
}
