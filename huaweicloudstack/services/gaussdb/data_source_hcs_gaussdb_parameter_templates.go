package gaussdb

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

// GaussDB GET /v3.1/{project_id}/configurations/{config_id}
func DataSourceConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConfigurationRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the parameter templates.`,
			},
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The parameter template ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the configuration.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the configuration.`,
			},
			"engine_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The engine version.`,
			},
			"instance_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance mode.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
			},
			"parameters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of configuration parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the parameter.`,
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The value of the parameter.`,
						},
						"need_restart": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the parameter needs to be restarted.`,
						},
						"readonly": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the parameter is read-only.`,
						},
						"value_range": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The value range of the parameter.`,
						},
						"data_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The data type of the parameter.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the parameter.`,
						},

						// The following parameters only supported from HCS 8.5.0 and later version
						"is_risk_parameter": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether this parameter is a risk parameter.`,
						},
						"risk_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of risk parameter.`,
						},
					},
				},
			},
		},
	}
}

func getConfiguration(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var httpUrl = "v3.1/{project_id}/configurations/{config_id}"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{config_id}", d.Get("template_id").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
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

	return respBody, nil
}

func flattenConfigurationParameters(parameters []interface{}) []map[string]interface{} {
	if len(parameters) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(parameters))
	for _, parameter := range parameters {
		result = append(result, map[string]interface{}{
			"name":              utils.PathSearch("name", parameter, nil),
			"value":             utils.PathSearch("value", parameter, nil),
			"need_restart":      utils.PathSearch("need_restart", parameter, nil),
			"readonly":          utils.PathSearch("readonly", parameter, nil),
			"value_range":       utils.PathSearch("value_range", parameter, nil),
			"data_type":         utils.PathSearch("data_type", parameter, nil),
			"description":       utils.PathSearch("description", parameter, nil),
			"is_risk_parameter": utils.PathSearch("is_risk_parameter", parameter, nil),
			"risk_description":  utils.PathSearch("risk_description", parameter, nil),
		})
	}

	return result
}

func dataSourceConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	configuration, err := getConfiguration(client, d)
	if err != nil {
		return diag.Errorf("error getting configuration: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", configuration, nil)),
		d.Set("description", utils.PathSearch("description", configuration, nil)),
		d.Set("engine_version", utils.PathSearch("engine_version", configuration, nil)),
		d.Set("instance_mode", utils.PathSearch("instance_mode", configuration, nil)),
		d.Set("created_at", utils.PathSearch("created_at", configuration, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", configuration, nil)),
		d.Set("parameters", flattenConfigurationParameters(
			utils.PathSearch("configuration_parameters", configuration, make([]interface{}, 0)).([]interface{}))))

	return diag.FromErr(mErr.ErrorOrNil())
}
