package authorizers

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const rootPath = "instances"

func rootURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(rootPath, instanceId, "authorizers")
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, id string) string {
	return c.ServiceURL(rootPath, instanceId, "authorizers", id)
}
