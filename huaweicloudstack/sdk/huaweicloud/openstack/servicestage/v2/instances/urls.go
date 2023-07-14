package instances

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const rootPath = "cas/applications"

func rootURL(client *golangsdk.ServiceClient, appId, componentId string) string {
	return client.ServiceURL(rootPath, appId, "components", componentId, "instances")
}

func resourceURL(client *golangsdk.ServiceClient, appId, componentId, instanceId string) string {
	return client.ServiceURL(rootPath, appId, "components", componentId, "instances", instanceId)
}

func actionURL(client *golangsdk.ServiceClient, appId, componentId, instanceId string) string {
	return client.ServiceURL(rootPath, appId, "components", componentId, "instances", instanceId, "action")
}
