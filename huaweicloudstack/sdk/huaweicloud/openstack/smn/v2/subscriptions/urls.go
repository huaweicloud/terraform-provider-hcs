package subscriptions

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func createURL(c *golangsdk.ServiceClient, topicUrn string) string {
	return c.ServiceURL("topics", topicUrn, "subscriptions")
}

func deleteURL(c *golangsdk.ServiceClient, subscriptionUrn string) string {
	return c.ServiceURL("subscriptions", subscriptionUrn)
}

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("subscriptions")
}

func listFromTopicURL(c *golangsdk.ServiceClient, topicUrn string) string {
	return c.ServiceURL("topics", topicUrn, "subscriptions")
}
