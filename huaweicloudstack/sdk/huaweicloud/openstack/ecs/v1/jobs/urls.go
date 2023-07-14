package jobs

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func rootURL(c *golangsdk.ServiceClient, jobId string) string {
	return c.ServiceURL("jobs", jobId)
}
