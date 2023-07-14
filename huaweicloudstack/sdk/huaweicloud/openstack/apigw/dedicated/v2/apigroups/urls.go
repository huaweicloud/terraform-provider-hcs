package apigroups

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const rootPath = "instances"

func rootURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(rootPath, instanceId, "api-groups")
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, groupId string) string {
	return c.ServiceURL(rootPath, instanceId, "api-groups", groupId)
}
