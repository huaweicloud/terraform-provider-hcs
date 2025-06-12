package vdc

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/vdc/v3/group"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
)

func DataSourceVdcGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVdcGroupRead,

		Schema: map[string]*schema.Schema{
			"vdc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"create_at": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceVdcGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	hcsConfig := config.GetHcsConfig(meta)
	vdcGroupClient, err := hcsConfig.VdcClient(hcsConfig.GetRegion(d))
	if err != nil {
		return diag.Errorf("error get vdc user group client: %s", err)
	}

	var groupInfo vdc.Group
	start := 0
	for {
		listOpts := vdc.ListReqParam{
			VdcID: d.Get("vdc_id").(string),
			Name:  d.Get("name").(string),
			Start: start,
			Limit: 100,
		}

		allGroups, total, err1 := vdc.GetGroupList(vdcGroupClient, listOpts).Extract()
		if err1 != nil {
			return diag.Errorf("error to retrieve vdc user groups: %s", err1)
		}

		// 是否有下一页数据
		hasNextPage := start+len(allGroups) < total

		// 精确查询数据
		filter := map[string]interface{}{
			"Name": d.Get("name").(string),
		}
		foundGroups, err2 := utils.FilterSliceWithZeroField(allGroups, filter)
		if err2 != nil {
			return diag.Errorf("filter vdc user group failed: %s", err2)
		}

		// 如果查到数据了
		if len(foundGroups) == 1 {
			groupInfo = foundGroups[0].(vdc.Group)
			d.SetId(groupInfo.ID)
			mErr := multierror.Append(nil,
				d.Set("name", groupInfo.Name),
				d.Set("domain_id", groupInfo.DomainId),
				d.Set("description", groupInfo.Description),
				d.Set("create_at", groupInfo.CreateAt),
			)
			if mErr.ErrorOrNil() != nil {
				return diag.Errorf("error setting vdc user group fields: %s", mErr)
			}
			break
		}

		// 如果没到最后一页，则继续查找
		if hasNextPage {
			start += listOpts.Limit
		} else { // 到了最后一页了
			if len(groupInfo.ID) < 1 {
				return diag.Errorf("not found vdc user group, please change you search conditions then try again!")
			}
			break
		}
	}

	return nil
}
