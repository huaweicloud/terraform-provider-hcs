package groups

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("vpc/security-groups")
}

func resourceURL(c *golangsdk.ServiceClient, secgroupId string) string {
	return c.ServiceURL("vpc/security-groups", secgroupId)
}
