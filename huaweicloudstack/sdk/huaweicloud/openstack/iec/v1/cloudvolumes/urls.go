package cloudvolumes

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("cloudvolumes")
}

func GetURL(c *golangsdk.ServiceClient, CloudVolumeID string) string {
	return c.ServiceURL("cloudvolumes", CloudVolumeID)
}

func ListVolumeTypeURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("cloudvolumes", "volume-types")
}
