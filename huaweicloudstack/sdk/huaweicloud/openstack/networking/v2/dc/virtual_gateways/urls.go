package virtual_gateways

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const (
	resourcePath = "virtual-gateways"
	rootpath     = "dcaas"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootpath, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootpath, resourcePath, id)
}
