package apis

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("apis")
}

func groupURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("apis", id)
}
