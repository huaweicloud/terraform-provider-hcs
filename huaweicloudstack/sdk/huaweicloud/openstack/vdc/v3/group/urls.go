package vdc

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

const rootPath = "rest/vdc/v3.2"

func GroupURL(c *golangsdk.ServiceClient, vdcId string) string {
	return c.ServiceURL(rootPath, "vdcs", vdcId, "groups")
}
