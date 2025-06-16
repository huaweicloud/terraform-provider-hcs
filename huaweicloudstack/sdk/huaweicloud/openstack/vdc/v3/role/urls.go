package role

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const vdcResourceBasePath = "rest/vdc/v3.0"

func GetVdcRoleURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(vdcResourceBasePath, "OS-ROLE/roles/third-party/roles")
}
