package user

import golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const resourcePath = "rest/vdc/v3.2"

func listOrCreateURL(c *golangsdk.ServiceClient, vdcId string) string {
	return c.ServiceURL(resourcePath, "vdcs", vdcId, "users")
}

func detailOrPutURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, "users", id)
}

func pwdURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("rest/vdc/v3.0/users", id, "reset-password")
}

func deleteURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL("rest/vdc/v3.0/users", id)
}
