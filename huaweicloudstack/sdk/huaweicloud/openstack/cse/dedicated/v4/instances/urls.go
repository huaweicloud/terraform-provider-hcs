package instances

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func rootURL(client *golangsdk.ServiceClient, serviceId string) string {
	return client.ServiceURL("registry", "microservices", serviceId, "instances")
}

func resourceURL(client *golangsdk.ServiceClient, serviceId, instanceId string) string {
	return client.ServiceURL("registry", "microservices", serviceId, "instances", instanceId)
}
