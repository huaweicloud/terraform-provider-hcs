package applications

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const rootPath = "cas/applications"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, appId string) string {
	return c.ServiceURL(rootPath, appId)
}

func configURL(c *golangsdk.ServiceClient, appId string) string {
	return c.ServiceURL(rootPath, appId, "configuration")
}
