package groups

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("groups")
}

func listUsersURL(client *golangsdk.ServiceClient, groupID string) string {
	return client.ServiceURL("groups", groupID, "users")
}

func getURL(client *golangsdk.ServiceClient, groupID string) string {
	return client.ServiceURL("groups", groupID)
}

func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("groups")
}

func updateURL(client *golangsdk.ServiceClient, groupID string) string {
	return client.ServiceURL("groups", groupID)
}

func deleteURL(client *golangsdk.ServiceClient, groupID string) string {
	return client.ServiceURL("groups", groupID)
}
