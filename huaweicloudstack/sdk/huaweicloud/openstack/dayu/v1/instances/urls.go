package instances

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const resourcePath = "instances"

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath, "onekey-purchase")
}
