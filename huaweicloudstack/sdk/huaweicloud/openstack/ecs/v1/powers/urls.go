package powers

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const rootURL = "cloudservers"

func actionURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootURL, "action")
}
