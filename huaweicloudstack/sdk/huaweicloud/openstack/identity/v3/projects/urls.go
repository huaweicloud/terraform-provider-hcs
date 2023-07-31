package projects

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("projects")
}

func getURL(client *golangsdk.ServiceClient, projectID string) string {
	return client.ServiceURL("projects", projectID)
}
