package instances

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(c.ProjectID, "instances")
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, "instances", id)
}

func resizeResourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, "instances", id, "resize")
}

func updatePasswordURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(c.ProjectID, "instances", id, "password")
}
