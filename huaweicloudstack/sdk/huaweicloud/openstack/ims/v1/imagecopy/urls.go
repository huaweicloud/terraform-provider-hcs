package imagecopy

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func withinRegionCopyURL(c *golangsdk.ServiceClient, imageId string) string {
	return c.ServiceURL("cloudimages", imageId, "copy")
}

func crossRegionCopyURL(c *golangsdk.ServiceClient, imageId string) string {
	return c.ServiceURL("cloudimages", imageId, "cross_region_copy")
}
