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
						"max_connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cps": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"qps": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
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
	curJson := utils.PathSearch("flavors", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil, nil
	}

	maxConnections := d.Get("max_connections").(int)
	cps := d.Get("cps").(int)
	qps := d.Get("qps").(int)
	bandwidth := d.Get("bandwidth").(int)

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	ids := make([]string, 0, len(curArray))
	for _, v := range curArray {
		rawConnection := utils.PathSearch("info.connection", v, float64(0)).(float64)
		if maxConnections > 0 && int(rawConnection) != maxConnections {
			continue
		}
		rawCps := utils.PathSearch("info.cps", v, float64(0)).(float64)
		if cps > 0 && int(rawCps) != cps {
			continue
		}
		rawQPS := utils.PathSearch("info.qps", v, float64(0)).(float64)
		if qps > 0 && int(rawQPS) != qps {
			continue
		}
		rawBandwidth := utils.PathSearch("info.bandwidth", v, float64(0)).(float64)
		if bandwidth > 0 && int(rawBandwidth) != bandwidth {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"id":              utils.PathSearch("id", v, nil),
			"name":            utils.PathSearch("name", v, nil),
			"type":            utils.PathSearch("type", v, nil),
			"shared":          utils.PathSearch("shared", v, nil),
			"flavor_sold_out": utils.PathSearch("flavor_sold_out", v, nil),
			"status":          utils.PathSearch("status", v, nil),
			"max_connections": rawConnection,
			"cps":             rawCps,
			"qps":             rawQPS,
			"bandwidth":       rawBandwidth,
		})
		ids = append(ids, utils.PathSearch("id", v, "").(string))
	}
	return rst, ids
}
