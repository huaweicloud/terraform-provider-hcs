package v1peering

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const resourcePath = "vpcpeering"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id)
}

func resourceListURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath) + "?peering_id=" + id
}

func acceptURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id, "accept")
}

func rejectURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, id, "reject")
}
