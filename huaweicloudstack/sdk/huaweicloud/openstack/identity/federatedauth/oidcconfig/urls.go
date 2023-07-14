package oidcconfig

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func resourceURL(c *golangsdk.ServiceClient, idpID string) string {
	return c.ServiceURL("v3.0", "OS-FEDERATION", "identity-providers", idpID, "openid-connect-config")
}
