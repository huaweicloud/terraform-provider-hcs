package snats

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("snat_rules")
}

func resourceURL(c *golangsdk.ServiceClient, ruleId string) string {
	return c.ServiceURL("snat_rules", ruleId)
}

func deleteURL(c *golangsdk.ServiceClient, gatewayId, ruleId string) string {
	return c.ServiceURL("nat_gateways", gatewayId, "snat_rules", ruleId)
}
