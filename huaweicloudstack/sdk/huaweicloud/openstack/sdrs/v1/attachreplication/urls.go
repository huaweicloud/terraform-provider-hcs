package attachreplication

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func createURL(c *golangsdk.ServiceClient, instanceID string) string {
	return c.ServiceURL("protected-instances", instanceID, "attachreplication")
}

func deleteURL(c *golangsdk.ServiceClient, instanceID string, replicationID string) string {
	return c.ServiceURL("protected-instances", instanceID, "detachreplication", replicationID)
}
