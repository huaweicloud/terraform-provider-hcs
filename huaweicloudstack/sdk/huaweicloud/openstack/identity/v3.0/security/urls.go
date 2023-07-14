package security

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const rootPath = "OS-SECURITYPOLICY"

func passwordPolicyURL(client *golangsdk.ServiceClient, domainID string) string {
	return client.ServiceURL(rootPath, "domains", domainID, "password-policy")
}
