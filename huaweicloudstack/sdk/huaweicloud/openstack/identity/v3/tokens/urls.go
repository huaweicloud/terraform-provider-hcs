package tokens

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func tokenURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("auth", "tokens")
}
