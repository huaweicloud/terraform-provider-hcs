package components

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const rootPath = "cas/applications"

func rootURL(c *golangsdk.ServiceClient, appId string) string {
	return c.ServiceURL(rootPath, appId, "components")
}

func resourceURL(c *golangsdk.ServiceClient, appId, componentId string) string {
	return c.ServiceURL(rootPath, appId, "components", componentId)
}
