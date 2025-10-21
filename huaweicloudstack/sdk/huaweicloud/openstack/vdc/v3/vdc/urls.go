package vdc

import golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const resourcePath = "rest/vdc/v3.0"

func listURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath, "vdcs")
}
