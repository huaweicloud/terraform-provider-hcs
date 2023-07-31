package users

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("users")
}

func getURL(client *golangsdk.ServiceClient, userID string) string {
	return client.ServiceURL("users", userID)
}

func listGroupsURL(client *golangsdk.ServiceClient, userID string) string {
	return client.ServiceURL("users", userID, "groups")
}

func listProjectsURL(client *golangsdk.ServiceClient, userID string) string {
	return client.ServiceURL("users", userID, "projects")
}
