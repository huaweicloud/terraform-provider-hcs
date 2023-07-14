package jobs

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

func rootURL(c *golangsdk.ServiceClient, clusterId string) string {
	return c.ServiceURL("clusters", clusterId, "job-executions")
}

func resourceURL(c *golangsdk.ServiceClient, clusterId, jobId string) string {
	return c.ServiceURL("clusters", clusterId, "job-executions", jobId)
}

func deleteURL(c *golangsdk.ServiceClient, clusterId string) string {
	return c.ServiceURL("clusters", clusterId, "job-executions", "batch-delete")
}
