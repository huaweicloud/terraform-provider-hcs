package instances

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const rootPath = "instances"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id)
}

func egressURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id, "nat-eip")
}

func ingressURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id, "eip")
}

func featureURL(c *golangsdk.ServiceClient, instanceId string) string {
	return c.ServiceURL(rootPath, instanceId, "features")
}
