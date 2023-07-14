package assignments

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("usg/dcs/corp/admin")
}

func resourceURL(client *golangsdk.ServiceClient, account string) string {
	return client.ServiceURL("usg/dcs/corp/admin", account)
}

func deleteURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("usg/dcs/corp/admin/delete")
}
