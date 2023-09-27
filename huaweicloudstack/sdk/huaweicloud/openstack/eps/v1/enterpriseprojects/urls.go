package enterpriseprojects

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const resourcePath = "enterprise-projects"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func migrateURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "resources-migrate")
}
