package sites

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

func ListURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL("sites")
}
