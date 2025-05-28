package flavors

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const (
	rootPath     = "elb"
	resourcePath = "flavors"
)

func listURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL(rootPath, resourcePath)
}

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(rootPath, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootPath, resourcePath, id)
}
