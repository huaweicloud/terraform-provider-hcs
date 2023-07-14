package images

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

// ListURL list iec image url
func ListURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL("images")
}
