package partitions

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const (
	rootPath     = "clusters"
	resourcePath = "partitions"
)

func rootURL(c *golangsdk.ServiceClient, clusterid string) string {
	return c.ServiceURL(rootPath, clusterid, resourcePath)
}

func resourceURL(c *golangsdk.ServiceClient, clusterid, partitionName string) string {
	return c.ServiceURL(rootPath, clusterid, resourcePath, partitionName)
}
