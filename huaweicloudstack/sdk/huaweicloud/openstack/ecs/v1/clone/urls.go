package clone

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const (
	rootPath     = "cloudservers"
	resourcePath = "action"
)

func cloneURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, id, resourcePath)
}

func jobURL(sc *golangsdk.ServiceClient, jobId string) string {
	return sc.ServiceURL("jobs", jobId)
}
