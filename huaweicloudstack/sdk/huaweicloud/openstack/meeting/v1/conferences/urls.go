package conferences

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const rootPath = "mmc/management"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath, "conferences")
}

func cycleURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath, "cycleconferences")
}

func subCycleURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath, "conferences", "cyclesubconf")
}

func showURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath, "conferences", "confDetail")
}

func onlineURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath, "conferences", "online", "confDetail")
}

func historyURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath, "conferences", "history", "confDetail")
}

func historiesURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath, "conferences", "history")
}

func controlURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath, "conference", "duration")
}
