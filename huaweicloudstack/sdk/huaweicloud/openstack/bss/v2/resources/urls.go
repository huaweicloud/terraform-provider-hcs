package resources

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func queryURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("orders/suscriptions/resources/query")
}

func autoRenewURL(c *golangsdk.ServiceClient, resourceId string) string {
	return c.ServiceURL("orders/subscriptions/resources/autorenew", resourceId)
}
