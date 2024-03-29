package volumes

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("volumes")
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("os-vendor-volumes", "detail")
}

func deleteURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("volumes", id)
}

func getURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("os-vendor-volumes", id)
}

func updateURL(c *golangsdk.ServiceClient, id string) string {
	return deleteURL(c, id)
}
