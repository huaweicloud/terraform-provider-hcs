package clouds

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("subscription/purchase/prepaid-cloud-waf")
}

func getURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("subscription")
}

func updateURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("subscription/batchalter/prepaid-cloud-waf")
}
