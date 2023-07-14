package backups

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const (
	cloudNativeRootPath = "cloudbackups"
	osNativeRootPath    = "backups"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(cloudNativeRootPath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(osNativeRootPath, id)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(osNativeRootPath, "detail")
}

func restoreURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(osNativeRootPath, id, "restore")
}
