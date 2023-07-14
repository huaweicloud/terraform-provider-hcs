package jobs

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func rootURL(client *golangsdk.ServiceClient, jobId string) string {
	return client.ServiceURL("cas/jobs", jobId)
}
