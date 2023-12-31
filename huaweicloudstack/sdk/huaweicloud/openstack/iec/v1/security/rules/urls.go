package rules

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

func CreateURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("security-group-rules")
}

func DeleteURL(c *golangsdk.ServiceClient, securityGroupRuleID string) string {
	return c.ServiceURL("security-group-rules", securityGroupRuleID)
}

func GetURL(c *golangsdk.ServiceClient, securityGroupRuleID string) string {
	return c.ServiceURL("security-group-rules", securityGroupRuleID)
}
