package instances

import golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const (
	rootPath = "instances"
)

// POST v1
func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

// GET v2, UPDATE v2
func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id)
}

// DELETE v1
func deleteURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("roma", rootPath, id)
}
