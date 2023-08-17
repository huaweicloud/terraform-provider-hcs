package eip

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/bandwidths"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/logp"
)

func DataSourceBandWidth() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBandWidthRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(5, 2000),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publicips": publicIPListComputedSchema(),
		},
	}
}

func dataSourceBandWidthRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := config.GetHcsConfig(meta)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloudStack VPC client: %s", err)
	}

	listOpts := bandwidths.ListOpts{
		EnterpriseProjectID: config.DataGetEnterpriseProjectID(d),
	}

	allBWs, err := bandwidths.List(vpcClient, listOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Unable to list HuaweiCloudStack bandwidths: %s", err)
	}
	if len(allBWs) == 0 {
		return fmtp.DiagErrorf("No HuaweiCloudStack bandwidth was found")
	}

	filter := map[string]interface{}{
		"Name": d.Get("name").(string),
	}
	if v, ok := d.GetOk("size"); ok {
		filter["Size"] = v
	}

	filterBWs, err := utils.FilterSliceWithField(allBWs, filter)
	if err != nil {
		return fmtp.DiagErrorf("filter bandwidths failed: %s", err)
	}
	if len(filterBWs) == 0 {
		return fmtp.DiagErrorf("No HuaweiCloudStack bandwidth was found by %+v", filter)
	}

	result := filterBWs[0].(bandwidths.BandWidth)
	logp.Printf("[DEBUG] Retrieved HuaweiCloudStack bandwidth %s: %+v", result.ID, result)
	d.SetId(result.ID)

	mErr := multierror.Append(
		d.Set("name", result.Name),
		d.Set("size", result.Size),
		d.Set("enterprise_project_id", result.EnterpriseProjectID),
		d.Set("share_type", result.ShareType),
		d.Set("bandwidth_type", result.BandwidthType),
		d.Set("status", result.Status),
		d.Set("publicips", flattenPublicIPs(result)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error setting bandwidth fields: %s", err)
	}

	return nil
}
