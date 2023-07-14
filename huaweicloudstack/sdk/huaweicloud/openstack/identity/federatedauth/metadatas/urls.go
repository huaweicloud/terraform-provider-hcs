package metadatas

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func resourceURL(c *golangsdk.ServiceClient, idpID string, protocolID string) string {
	return c.ServiceURL("v3-ext", "OS-FEDERATION", "identity_providers", idpID, "protocols", protocolID, "metadata")
}
