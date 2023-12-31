package namespaces

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("manage", "namespaces")
}

func resourceURL(client *golangsdk.ServiceClient, name string) string {
	return client.ServiceURL("manage", "namespaces", name)
}

func accessURL(client *golangsdk.ServiceClient, namespace string) string {
	return client.ServiceURL("manage", "namespaces", namespace, "access")
}
