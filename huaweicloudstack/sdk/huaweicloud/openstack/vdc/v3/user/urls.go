package user

import golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const resourcePath = "rest/vdc/v3.2/vdcs"

func rootURL(c *golangsdk.ServiceClient, vdcId string) string {
	return c.ServiceURL(resourcePath, vdcId, "users")
}
