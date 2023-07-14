package providers

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("v3", "OS-FEDERATION", "identity_providers", id)
}
