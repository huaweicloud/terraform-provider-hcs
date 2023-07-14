package propagations

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func enableURL(client *golangsdk.ServiceClient, instanceId, routeTableId string) string {
	return client.ServiceURL("enterprise-router", instanceId, "route-tables", routeTableId, "enable-propagations")
}

func queryURL(client *golangsdk.ServiceClient, instanceId, routeTableId string) string {
	return client.ServiceURL("enterprise-router", instanceId, "route-tables", routeTableId, "propagations")
}

func disableURL(client *golangsdk.ServiceClient, instanceId, routeTableId string) string {
	return client.ServiceURL("enterprise-router", instanceId, "route-tables", routeTableId, "disable-propagations")
}
