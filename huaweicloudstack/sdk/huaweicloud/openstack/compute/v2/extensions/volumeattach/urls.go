package volumeattach

import (
	"strings"

	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

const resourcePath = "os-volume_attachments"

func resourceURL(c *golangsdk.ServiceClient, serverID string) string {
	return c.ServiceURL("servers", serverID, resourcePath)
}

func listURL(c *golangsdk.ServiceClient, serverID string) string {
	return resourceURL(c, serverID)
}

func createURL(c *golangsdk.ServiceClient, serverID string) string {
	return resourceURL(c, serverID)
}

func getURL(c *golangsdk.ServiceClient, serverID, aID string) string {
	return c.ServiceURL("servers", serverID, resourcePath, aID)
}

func deleteURL(c *golangsdk.ServiceClient, serverID, aID string) string {
	url := strings.Replace(getURL(c, serverID, aID), "/v2.1/", "/v2/", 1)
	return url
}
