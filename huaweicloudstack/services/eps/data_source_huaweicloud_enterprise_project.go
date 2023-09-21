package eps

import (
	"context"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/eps/v1/enterpriseprojects"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
)

func DataSourceEnterpriseProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnterpriseProjectRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vdc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"inherit": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "prod",
			},
			"status": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"query_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "list",
			},
			"auth_action": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contain_default": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"offset": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
			},
			"limit": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "1000",
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "created_at",
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "asc",
			},
			"instances": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delete_flag": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vdc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vdc_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceEnterpriseProjectRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := config.GetHcsConfig(meta)
	region := config.GetRegion(d)
	epsClient, err := config.EnterpriseProjectClient(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud eps client %s", err)
	}

	listOpts := enterpriseprojects.ListOpts{
		Name:           d.Get("name").(string),
		ID:             d.Get("id").(string),
		IDs:            d.Get("ids").(string),
		DomainId:       d.Get("domain_id").(string),
		VdcId:          d.Get("vdc_id").(string),
		Inherit:        d.Get("inherit").(bool),
		ProjectId:      d.Get("project_id").(string),
		Type:           d.Get("type").(string),
		Status:         d.Get("status").(int),
		QueryType:      d.Get("query_type").(string),
		AuthAction:     d.Get("auth_action").(string),
		ContainDefault: d.Get("contain_default").(bool),
		Offset:         d.Get("offset").(string),
		Limit:          d.Get("limit").(string),
		SortKey:        d.Get("sort_key").(string),
		SortDir:        d.Get("sort_dir").(string),
	}
	projects, err := enterpriseprojects.List(epsClient, listOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error retrieving enterprise projects %s", err)
	}
	ids := filterIDs(projects)
	// Save the data source ID using a hash code constructed using all instance IDs.
	d.SetId(hashcode.Strings(ids))
	result := make([]map[string]interface{}, len(projects))
	for i, item := range projects {
		server := map[string]interface{}{
			"id":           item.ID,
			"name":         item.Name,
			"description":  item.Description,
			"type":         item.Type,
			"delete_flag":  item.DeleteFlag,
			"status":       item.Status,
			"created_at":   item.CreatedAt,
			"updated_at":   item.UpdatedAt,
			"domain_id":    item.DomainId,
			"vdc_id":       item.VdcId,
			"project_id":   item.ProjectId,
			"project_name": item.ProjectName,
		}
		result[i] = server
	}
	mErr := multierror.Append(nil,
		d.Set("instances", result),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func filterIDs(projects []enterpriseprojects.Project) []string {
	ids := make([]string, 0, len(projects))

	for _, item := range projects {
		ids = append(ids, item.ID)
	}

	return ids
}
