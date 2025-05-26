package vdc

import (
	"context"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vdc/v3/user"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
)

func DataResourceVdcUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataResourceVdcUserRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: false,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Required: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"vdc_id": {
				Type:     schema.TypeString,
				Optional: false,
				Required: true,
			},
			"ldap_id": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"create_at": {
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
			},
			"auth_type": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"access_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"top_vdc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},
			"start": {
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
				Default:  0,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Required: false,
				Default:  100,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Default:  "name",
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
				Default:  "asc",
			},
		},
	}
}

func dataResourceVdcUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := config.GetHcsConfig(meta)
	region := config.GetRegion(d)
	vdcUserClient, err := config.VdcClient(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud vdc user client %s", err)
	}

	vdcId := d.Get("vdc_id").(string)
	searchName := d.Get("name").(string)
	listOpts := user.ListOpts{
		Name:    searchName,
		Start:   user.Start,
		Limit:   user.Limit,
		SortKey: user.SortKey,
		SortDir: user.SortDir,
	}
	res, err := user.List(vdcUserClient, vdcId, listOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error retrieving vdc user %s", err)
	}

	if len(res.Users) == 0 {
		return diag.Errorf("your query returned no results, please change your search criteria and try again.")
	} else {
		allUsers := res.Users
		for {
			result := findSearchName(res.Users, searchName)

			if result.ID != "" {
				// Save the data source ID using a hash code constructed using all instance IDs.
				d.SetId(result.ID)
				d.Set("domain_id", result.DomainId)
				d.Set("name", result.Name)
				d.Set("display_name", result.DisplayName)
				d.Set("enabled", result.Enabled)
				d.Set("description", result.Description)
				d.Set("vdc_id", result.VdcId)
				d.Set("ldap_id", result.LdapId)
				d.Set("create_at", result.CreateAt)
				typeVal := user.AuthType[result.AuthType]
				d.Set("auth_type", typeVal)
				modeVal := user.AccessMode[result.AccessMode]
				d.Set("access_mode", modeVal)
				d.Set("top_vdc_id", result.TopVdcId)

				return diag.FromErr(nil)
			} else {
				if len(allUsers) < int(res.Total) {
					listOpts.Start = listOpts.Start + 1
					res, err = user.List(vdcUserClient, vdcId, listOpts).Extract()
					if err != nil {
						return fmtp.DiagErrorf("Error retrieving vdc user %s", err)
					}
					allUsers = append(allUsers, res.Users...)
				} else {
					return diag.Errorf("您的查询返回了多个结果，请更改搜索名称进行精确查询。")
				}
			}
		}

	}
}

func findSearchName(users []user.VdcUserModel, searchName string) user.VdcUserModel {
	var userModel user.VdcUserModel

	for _, val := range users {
		if val.Name == searchName {
			userModel = val
			break
		}
	}

	return userModel
}
