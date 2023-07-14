package keypairs

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

func rootURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("os-keypairs")
}

func DeleteURL(c *golangsdk.ServiceClient, KeyPairName string) string {
	return c.ServiceURL("os-keypairs", KeyPairName)
}

func GetURL(c *golangsdk.ServiceClient, KeyPairName string) string {
	return c.ServiceURL("os-keypairs", KeyPairName)
}
