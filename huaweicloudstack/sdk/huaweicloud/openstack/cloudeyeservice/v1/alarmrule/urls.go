package alarmrule

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const (
	rootPath = "alarms"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id)
}

func actionURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id, "action")
}
