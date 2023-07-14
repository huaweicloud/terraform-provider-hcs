package persistentvolumeclaims

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/cce/v3/addons"
)

const rootPath = "namespaces"

func createURL(client *golangsdk.ServiceClient, clusterId, ns string) string {
	return addons.CCEServiceURL(client, clusterId, rootPath, ns, "persistentvolumeclaims")
}

func listURL(client *golangsdk.ServiceClient, clusterId, ns string) string {
	return addons.CCEServiceURL(client, clusterId, rootPath, ns, "persistentvolumeclaims")
}

func deleteURL(client *golangsdk.ServiceClient, clusterId, ns, name string) string {
	return addons.CCEServiceURL(client, clusterId, rootPath, ns, "persistentvolumeclaims", name)
}
