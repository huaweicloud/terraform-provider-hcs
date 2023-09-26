package vpc

import (
	"context"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/common"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	v1rules "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/security/rules"
	v1groups "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v1/security/securitygroups"
	v3groups "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/networking/v3/security/groups"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/utils/logp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceNetworkingSecGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingSecGroupRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"secgroup_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rules": securityGroupRuleSchema,
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getRuleListByGroupId(client *golangsdk.ServiceClient, groupId string) ([]map[string]interface{}, error) {
	listOpts := v1rules.ListOpts{
		SecurityGroupId: groupId,
	}
	resp, err := v1rules.List(client, listOpts)
	if err != nil {
		return nil, err
	}
	return flattenSecurityGroupRuleV1(resp)
}

func flattenSecurityGroupRuleV1(rules []v1rules.SecurityGroupRule) ([]map[string]interface{}, error) {
	sgRules := make([]map[string]interface{}, len(rules))
	for i, rule := range rules {
		ruleInfo := map[string]interface{}{
			"id":               rule.ID,
			"direction":        rule.Direction,
			"protocol":         rule.Protocol,
			"ethertype":        rule.Ethertype,
			"remote_ip_prefix": rule.RemoteIpPrefix,
			"remote_group_id":  rule.RemoteGroupId,
			"description":      rule.Description,
		}

		sgRules[i] = ruleInfo
	}

	return sgRules, nil
}

func dataSourceNetworkingSecGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := common.GetRegion(d, cfg)
	v3Client, err := cfg.NetworkingV3Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking v3 client: %s", err)
	}

	listOpts := v3groups.ListOpts{
		ID:   d.Get("secgroup_id").(string),
		Name: d.Get("name").(string),
	}

	allSecGroups, err := v3groups.List(v3Client, listOpts)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			// If the v3 API does not exist or has not been published in the specified region, set again using v1 API.
			return dataSourceNetworkingSecGroupReadV1(ctx, d, meta)
		}
		return fmtp.DiagErrorf("Unable to get security groups list: %s", err)
	}

	if len(allSecGroups) < 1 {
		return fmtp.DiagErrorf("No Security Group found.")
	}

	if len(allSecGroups) > 1 {
		return fmtp.DiagErrorf("More than one Security Groups found.")
	}

	secGroup := allSecGroups[0]
	d.SetId(secGroup.ID)
	logp.Printf("[DEBUG] Retrieved Security Group (%s) by v3 client: %v", d.Id(), secGroup)

	v1Client, err := cfg.NetworkingV1Client(region)
	rules, err := getRuleListByGroupId(v1Client, secGroup.ID)
	if err != nil {
		return diag.FromErr(err)
	}
	logp.Printf("[DEBUG] The rules list of security group (%s) is: %v", d.Id(), rules)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", secGroup.Name),
		d.Set("description", secGroup.Description),
		d.Set("rules", rules),
		d.Set("created_at", secGroup.CreatedAt),
		d.Set("updated_at", secGroup.UpdatedAt),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func dataSourceNetworkingSecGroupReadV1(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := config.GetHcsConfig(meta)
	region := common.GetRegion(d, cfg)
	v1Client, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud networking v1 client: %s", err)
	}

	listOpts := v1groups.ListOpts{
		EnterpriseProjectId: cfg.DataGetEnterpriseProjectID(d),
	}

	pages, err := v1groups.List(v1Client, listOpts).AllPages()
	if err != nil {
		return diag.FromErr(err)
	}
	allGroups, err := v1groups.ExtractSecurityGroups(pages)
	if err != nil {
		return fmtp.DiagErrorf("Error retrieving security groups list: %s", err)
	}
	if len(allGroups) == 0 {
		return fmtp.DiagErrorf("No sucurity group found, please change your search criteria and try again.")
	}
	logp.Printf("[DEBUG] The retrieved group list is: %v", allGroups)

	filter := map[string]interface{}{
		"ID":   d.Get("secgroup_id"),
		"Name": d.Get("name"),
	}
	filterGroups, err := utils.FilterSliceWithField(allGroups, filter)
	if err != nil {
		return fmtp.DiagErrorf("Erroring filting security groups list: %s", err)
	}
	if len(filterGroups) < 1 {
		return fmtp.DiagErrorf("No Security Group found.")
	}
	if len(filterGroups) > 1 {
		return fmtp.DiagErrorf("More than one Security Groups found.")
	}

	resp := filterGroups[0].(v1groups.SecurityGroup)
	d.SetId(resp.ID)

	rules := flattenSecurityGroupRulesV1(&resp)
	logp.Printf("[DEBUG] The retrieved rules list is: %v", rules)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
		d.Set("rules", rules),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
