package mappings

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("v3", "OS-FEDERATION", "mappings", id)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("v3", "OS-FEDERATION", "mappings")
}
