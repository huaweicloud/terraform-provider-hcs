package flavors

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("flavors")
}
