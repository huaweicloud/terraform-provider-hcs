package dc_endpoint_groups

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const (
	resourcePath = "dc-endpoint-groups"
	rootpath     = "dcaas"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootpath, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootpath, resourcePath, id)
}
