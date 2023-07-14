package interfaces

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("dcaas/virtual-interfaces")
}

func resourceURL(client *golangsdk.ServiceClient, interfaceId string) string {
	return client.ServiceURL("dcaas/virtual-interfaces", interfaceId)
}
