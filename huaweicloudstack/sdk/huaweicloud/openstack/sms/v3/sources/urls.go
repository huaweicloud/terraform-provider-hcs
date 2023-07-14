package sources

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("sources")
}

func getURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("sources", id)
}
