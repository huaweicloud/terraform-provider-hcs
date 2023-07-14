package persistentvolumeclaims

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const rootPath = "namespaces"

func rootURL(client *golangsdk.ServiceClient, ns string) string {
	return client.ServiceURL(rootPath, ns, "extended-persistentvolumeclaims")
}

func resourceURL(client *golangsdk.ServiceClient, ns, name string) string {
	return client.ServiceURL(rootPath, ns, "persistentvolumeclaims", name)
}
