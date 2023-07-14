package plugins

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func rootURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL("instances", instanceId, "plugins")
}

func resourceURL(c *golangsdk.ServiceClient, instanceId, pluginId string) string {
	return c.ServiceURL("instances", instanceId, "plugins", pluginId)
}

func bindURL(c *golangsdk.ServiceClient, instanceId, pluginId string) string {
	return c.ServiceURL("instances", instanceId, "plugins", pluginId, "attach")
}

func listBindURL(c *golangsdk.ServiceClient, instanceId, pluginId string) string {
	return c.ServiceURL("instances", instanceId, "plugins", pluginId, "attached-apis")
}

func unbindURL(c *golangsdk.ServiceClient, instanceId, pluginId string) string {
	return c.ServiceURL("instances", instanceId, "plugins", pluginId, "detach")
}
