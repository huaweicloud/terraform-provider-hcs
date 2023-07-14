package subnets

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("subnets")
}

func DeleteURL(c *golangsdk.ServiceClient, subnetId string) string {
	return c.ServiceURL("subnets", subnetId)
}

func GetURL(c *golangsdk.ServiceClient, subnetId string) string {
	return c.ServiceURL("subnets", subnetId)
}

func UpdateURL(c *golangsdk.ServiceClient, subnetId string) string {
	return c.ServiceURL("subnets", subnetId)
}
