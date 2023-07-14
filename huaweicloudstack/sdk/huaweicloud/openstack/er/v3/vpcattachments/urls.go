package vpcattachments

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func rootURL(client *golangsdk.ServiceClient, instanceId string) string {
	return client.ServiceURL("enterprise-router", instanceId, "vpc-attachments")
}

func resourceURL(client *golangsdk.ServiceClient, instanceId, attachmentId string) string {
	return client.ServiceURL("enterprise-router", instanceId, "vpc-attachments", attachmentId)
}
