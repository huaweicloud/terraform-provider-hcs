package securities

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func rootURL(c *golangsdk.ServiceClient, instanceId, path string) string {
	return c.ServiceURL("instances", instanceId, path)
}
