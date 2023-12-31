package namespaces

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const rootPath = "namespaces"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath)
}

func resourceURL(client *golangsdk.ServiceClient, name string) string {
	return client.ServiceURL(rootPath, name)
}
