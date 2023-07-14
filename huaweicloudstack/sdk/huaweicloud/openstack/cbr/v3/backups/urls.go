package backups

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func resourceURL(c *golangsdk.ServiceClient, backupId string) string {
	return c.ServiceURL("backups", backupId)
}
