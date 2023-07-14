package flavors

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func listURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL("cloudservers", "flavors")
}
