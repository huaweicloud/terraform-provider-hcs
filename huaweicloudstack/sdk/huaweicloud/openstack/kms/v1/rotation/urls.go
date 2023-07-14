package rotation

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const (
	resourcePath = "kms"
)

func enableURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath, "enable-key-rotation")
}

func disableURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath, "disable-key-rotation")
}

func intervalURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath, "update-key-rotation-interval")
}

func getURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath, "get-key-rotation-status")
}
