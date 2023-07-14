package groups

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("security-groups")
}

func DeleteURL(c *golangsdk.ServiceClient, securityGroupID string) string {
	return c.ServiceURL("security-groups", securityGroupID)
}

func GetURL(c *golangsdk.ServiceClient, securityGroupID string) string {
	return c.ServiceURL("security-groups", securityGroupID)
}
