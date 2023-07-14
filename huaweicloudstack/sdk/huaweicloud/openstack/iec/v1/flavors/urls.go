package flavors

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

func ListURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("cloudservers", "flavors")
}
