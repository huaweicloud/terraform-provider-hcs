package securitygroups

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

func CreateURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, "security-groups")
}

func DeleteURL(c *golangsdk.ServiceClient, securityGroupId string) string {
	return c.ServiceURL(c.ProjectID, "security-groups", securityGroupId)
}

func GetURL(c *golangsdk.ServiceClient, securityGroupId string) string {
	return c.ServiceURL(c.ProjectID, "security-groups", securityGroupId)
}

func ListURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, "security-groups")
}
