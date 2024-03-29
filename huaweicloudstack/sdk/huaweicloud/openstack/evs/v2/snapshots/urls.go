package snapshots

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("snapshots")
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("os-vendor-snapshots/detail")
}

func deleteURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("snapshots", id)
}

func updateURL(c *golangsdk.ServiceClient, id string) string {
	return deleteURL(c, id)
}

func getURL(c *golangsdk.ServiceClient, id string) string {
	return deleteURL(c, id)
}
