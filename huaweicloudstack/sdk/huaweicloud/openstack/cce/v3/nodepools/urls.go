package nodepools

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const (
	rootPath     = "clusters"
	resourcePath = "nodepools"
)

func rootURL(c *golangsdk.ServiceClient, clusterid string) string {
	return c.ServiceURL(rootPath, clusterid, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, clusterid, nodepoolid string) string {
	return c.ServiceURL(rootPath, clusterid, resourcePath, nodepoolid)
}
