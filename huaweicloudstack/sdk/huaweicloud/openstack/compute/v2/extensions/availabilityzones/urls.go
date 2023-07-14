package availabilityzones

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("os-availability-zone")
}

func listDetailURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("os-availability-zone", "detail")
}
