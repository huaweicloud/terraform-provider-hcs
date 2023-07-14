package queues

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const (
	resourcePath = "queues"
	actionPath   = "action"
)

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, queueName string) string {
	return c.ServiceURL(resourcePath, queueName)
}

func queryAllURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func ActionURL(c *golangsdk.ServiceClient, queueName string) string {
	return c.ServiceURL(resourcePath, queueName, actionPath)
}
