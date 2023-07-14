package sqlfilter

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func updateURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "sql-filter", "switch")
}

func getURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "sql-filter", "switch")
}
