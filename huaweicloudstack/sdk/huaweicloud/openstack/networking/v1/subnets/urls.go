package subnets

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const (
	resourcePath = "subnets"
	rootpath     = "vpcs"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id)
}

func updateURL(c *golangsdk.ServiceClient, vpcid, id string) string {
	return c.ServiceURL(c.ProjectID, rootpath, vpcid, resourcePath, id)
}
