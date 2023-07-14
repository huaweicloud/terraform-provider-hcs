package metadata

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const rootPath = "cas/metadata"

func runtimeURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, "runtimes")
}

func flavorURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, "flavors")
}
