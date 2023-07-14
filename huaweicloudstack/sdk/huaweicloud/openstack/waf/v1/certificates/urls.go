package certificates

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("certificate")
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("certificate", id)
}
