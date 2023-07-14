package auth

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const rootPath = "usg/acs/auth/appauth"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath)
}
