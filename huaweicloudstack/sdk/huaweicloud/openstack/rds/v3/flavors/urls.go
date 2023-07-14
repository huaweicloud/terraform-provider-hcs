package flavors

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func listURL(sc *golangsdk.ServiceClient, databasename string) string {
	return sc.ServiceURL("flavors", databasename)
}
