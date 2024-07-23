package vpcs

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const addCidrPath = "add-extend-cidr"

const removeCidrPath = "remove-extend-cidr"

const resourcePath = "vpc/vpcs"

func addCidrURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, addCidrPath)
}

func removeCidrURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, removeCidrPath)
}

func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}
