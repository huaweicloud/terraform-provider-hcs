package domains

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func getURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("active-domains")
}
