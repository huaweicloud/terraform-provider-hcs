package gateways

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("private-nat/gateways")
}

func resourceURL(c *golangsdk.ServiceClient, ruleId string) string {
	return c.ServiceURL("private-nat/gateways", ruleId)
}
