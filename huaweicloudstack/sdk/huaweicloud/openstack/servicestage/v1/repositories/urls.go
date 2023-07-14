package repositories

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("git/auths")
}

func resourceURL(c *golangsdk.ServiceClient, name string) string {
	return c.ServiceURL("git/auths", name)
}

func passwordURL(c *golangsdk.ServiceClient, rType string) string {
	return c.ServiceURL("git/auths", rType, "password")
}

func personalURL(c *golangsdk.ServiceClient, rType string) string {
	return c.ServiceURL("git/auths", rType, "personal")
}
