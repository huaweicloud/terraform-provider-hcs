package v1peeringroute

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const (
	resourcePath = "vpcpeering"
	router       = "router"
	remove       = "removeroutes"
	query        = "queryroutes"
	add          = "addroutes"
)

func queryURL(c *golangsdk.ServiceClient, vpcI string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, router, vpcI, query)
}

func addURL(c *golangsdk.ServiceClient, vpcId string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, router, vpcId, add)
}

func removeURL(c *golangsdk.ServiceClient, vpcId string) string {
	return c.ServiceURL(c.ProjectID, resourcePath, router, vpcId, remove)
}
