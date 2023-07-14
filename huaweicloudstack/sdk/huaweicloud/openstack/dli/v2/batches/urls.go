package batches

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const rootPath = "batches"

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath)
}

func resourceURL(c *golangsdk.ServiceClient, batchId string) string {
	return c.ServiceURL(rootPath, batchId)
}

func stateURL(c *golangsdk.ServiceClient, batchId string) string {
	return c.ServiceURL(rootPath, batchId, "state")
}

func logURL(c *golangsdk.ServiceClient, batchId string) string {
	return c.ServiceURL(rootPath, batchId, "log")
}
