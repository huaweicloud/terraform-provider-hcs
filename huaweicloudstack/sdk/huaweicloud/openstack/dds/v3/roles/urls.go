package roles

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func rootURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "db-role")
}

func resourceURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "db-roles")
}
