package dnats

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const rootPath = "private-nat/dnat-rules"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, ruleId string) string {
	return c.ServiceURL(rootPath, ruleId)
}
