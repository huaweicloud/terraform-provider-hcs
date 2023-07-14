package users

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("users")
}

func resourceURL(c *golangsdk.ServiceClient, userId string) string {
	return c.ServiceURL("users", userId)
}
