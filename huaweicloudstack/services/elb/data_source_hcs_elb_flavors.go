package elb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/pagination"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

// @API ELB GET /v3/{project_id}/elb/flavors
func DataSourceElbFlavorsV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceElbFlavorsV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"flavor_sold_out": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_connections": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cps": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"qps": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// Computed values.
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"shared": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"flavor_sold_out": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"info": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceElbFlavorsV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	var (
		httpUrl = "v3/{project_id}/elb/flavors"
		product = "vpc"
	)
	listFlavorClient, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	listPath := listFlavorClient.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", listFlavorClient.ProjectID)
	listQueryParams := buildListFlavorsQueryParams(d)
	listPath += listQueryParams
	listResp, err := pagination.ListAllItems(
		listFlavorClient,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving ELB flavors: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	flavors, ids := flattenListFlavorsBody(d, listRespBody)
	mErr := multierror.Append(
		d.Set("region", cfg.GetRegion(d)),
		d.Set("ids", ids),
		d.Set("flavors", flavors),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListFlavorsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("flavor_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOkExists("shared"); ok {
		shared := v.(bool)
		res = fmt.Sprintf("%s&shared=%v", res, shared)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenListFlavorsBody(d *schema.ResourceData, resp interface{}) ([]interface{}, []string) {
	if resp == nil {
		return nil, nil
	}
	curJson := utils.PathSearch("flavors", resp, []interface{}{})
	if curJson == nil {
		return nil, nil
	}
	// represent info filters condition when querying
	hasMaxConnections := false
	hasCps := false
	hasQps := false
	hasBandwidth := false

	maxConnections := 0
	cps := 0
	qps := 0
	bandwidth := 0

	if v, ok := d.GetOkExists("max_connections"); ok {
		maxConnections = v.(int)
		hasMaxConnections = true
	}
	if v, ok := d.GetOkExists("cps"); ok {
		cps = v.(int)
		hasCps = true
	}
	if v, ok := d.GetOkExists("qps"); ok {
		qps = v.(int)
		hasQps = true
	}
	if v, ok := d.GetOkExists("bandwidth"); ok {
		bandwidth = v.(int)
		hasBandwidth = true
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	ids := make([]string, 0, len(curArray))

	for _, v := range curArray {
		flavorMap := map[string]interface{}{
			"id":              utils.PathSearch("id", v, nil),
			"name":            utils.PathSearch("name", v, nil),
			"type":            utils.PathSearch("type", v, nil),
			"shared":          utils.PathSearch("shared", v, nil),
			"flavor_sold_out": utils.PathSearch("flavor_sold_out", v, nil),
			"status":          utils.PathSearch("status", v, nil),
		}

		// Flatten the info field and include it regardless of filtering
		rawInfo := utils.PathSearch("info", v, nil)
		info := map[string]interface{}{}
		if rawInfoMap, ok := rawInfo.(map[string]interface{}); ok && rawInfoMap != nil {
			for k, vv := range rawInfoMap {
				info[k] = vv
			}
		}
		// filter qos key
		if hasMaxConnections {
			if val, ok := info["connection"].(float64); ok {
				if int(val) != maxConnections {
					continue
				}
			} else {
				continue
			}
		}
		if hasCps {
			if val, ok := info["cps"].(float64); ok {
				if int(val) != cps {
					continue
				}
			} else {
				continue
			}
		}
		if hasQps {
			if val, ok := info["qps"].(float64); ok {
				if int(val) != qps {
					continue
				}
			} else {
				continue
			}
		}
		if hasBandwidth {
			if val, ok := info["bandwidth"].(float64); ok {
				if int(val) != bandwidth {
					continue
				}
			} else {
				continue
			}
		}

		flavorMap["info"] = info
		rst = append(rst, flavorMap)
		ids = append(ids, utils.PathSearch("id", v, "").(string))
	}
	return rst, ids
}
