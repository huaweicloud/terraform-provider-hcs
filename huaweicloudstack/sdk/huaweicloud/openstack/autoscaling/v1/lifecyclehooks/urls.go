package lifecyclehooks

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const rootPath = "scaling_lifecycle_hook"

func rootURL(client *golangsdk.ServiceClient, groupID string) string {
	return client.ServiceURL(rootPath, groupID)
}

func resourceURL(client *golangsdk.ServiceClient, groupID, hookName string) string {
	return client.ServiceURL(rootPath, groupID, hookName)
}

func listURL(client *golangsdk.ServiceClient, groupID string) string {
	return client.ServiceURL(rootPath, groupID, "list")
}
