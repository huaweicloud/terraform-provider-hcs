package configurations

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func applyURL(c *golangsdk.ServiceClient, configId string) string {
	return c.ServiceURL("configurations", configId, "apply")
}

func getURL(c *golangsdk.ServiceClient, configId string) string {
	return c.ServiceURL("configurations", configId)
}

func instanceConfigURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "configurations")
}
