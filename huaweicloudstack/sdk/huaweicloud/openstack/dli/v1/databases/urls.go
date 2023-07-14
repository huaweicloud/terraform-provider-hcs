package databases

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const rootPath = "databases"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, dbName string) string {
	return c.ServiceURL(rootPath, dbName)
}

func userURL(c *golangsdk.ServiceClient, dbName string) string {
	return c.ServiceURL(rootPath, dbName, "owner")
}
