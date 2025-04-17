package cfw

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

// @API CFW POST /v1/{project_id}/acl-rule/count
func DataSourceCfwProtectionRuleHitCount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCfwProtectionRuleHitCountRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rule_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"fw_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			// Attribute
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_hit_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCfwProtectionRuleHitCountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := cfg.GetRegion(d)
	var mErr *multierror.Error

	// get cfw Protection Rule Hit Count
	var (
		getProtectionRuleHitCountUrl     = "v1/{project_id}/acl-rule/count"
		getProtectionRuleHitCountProduct = "cfw"
	)
	client, err := cfg.NewServiceClient(getProtectionRuleHitCountProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	// build request
	getProtectionRuleHitCountPath := client.Endpoint + getProtectionRuleHitCountUrl
	getProtectionRuleHitCountPath = strings.ReplaceAll(getProtectionRuleHitCountPath, "{project_id}", client.ProjectID)
	getProtectionRuleHitCountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildRuleHitCountBodyParams(d.Get("rule_ids").([]interface{})),
	}

	// get result
	var res []interface{}
	offset := 0
	for {
		getPath := getProtectionRuleHitCountPath + buildGetHitCountQueryParams(d, offset)
		getResp, err := client.Request("POST", getPath, &getProtectionRuleHitCountOpt)
		if err != nil {
			return diag.Errorf("error request CFW Protection rule hit count: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		ruleCountList := flattenRuleHitCount(getRespBody.(map[string]interface{}))
		res = append(res, ruleCountList...)
		totalAccount := utils.PathSearch("total", getRespBody, float64(0)).(float64)
		if len(res) >= int(totalAccount) {
			break
		}
		offset++
	}

	// set parameters
	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID of data source hcs_cfw_protection_rule_hit_count: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("records", res),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetHitCountQueryParams(d *schema.ResourceData, offset int) string {
	queryParam := "limit=100"
	if fwInstanceId, ok := d.GetOk("fw_instance_id"); ok {
		queryParam = fmt.Sprintf("%s&fw_instance_id=%s", queryParam, fwInstanceId.(string))
	}
	if epsId, ok := d.GetOk("enterprise_project_id"); ok {
		queryParam = fmt.Sprintf("%s&enterprise_project_id=%s", queryParam, epsId.(string))
	} else {
		queryParam = fmt.Sprintf("%s&enterprise_project_id=0", queryParam)
	}

	return fmt.Sprintf("?offset=%d&%s", offset, queryParam)
}

func buildRuleHitCountBodyParams(ruleIds []interface{}) map[string]interface{} {
	var ids []string
	for _, id := range ruleIds {
		if str, ok := id.(string); ok {
			ids = append(ids, str)
		}
	}

	return map[string]interface{}{
		"rule_ids": ids,
	}
}

func flattenRuleHitCount(resp map[string]interface{}) []interface{} {
	var records []interface{}
	recordsList := resp["data"].(map[string]interface{})["records"].([]interface{})

	for _, rec := range recordsList {
		recMap := rec.(map[string]interface{})

		ruleId, _ := recMap["rule_id"].(string)
		ruleHitCount := int(recMap["rule_hit_count"].(float64))

		records = append(records, map[string]interface{}{
			"rule_id":        ruleId,
			"rule_hit_count": ruleHitCount,
		})
	}

	return records
}
