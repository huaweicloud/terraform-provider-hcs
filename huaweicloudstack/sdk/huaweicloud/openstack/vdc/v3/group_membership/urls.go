package group_membership

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	url2 "net/url"
)

const rootPath = "rest/vdc/v3.2"

func GroupMemberShipURL(c *golangsdk.ServiceClient, groupId string) string {
	return c.ServiceURL(rootPath, "groups", url2.PathEscape(groupId), "users")
}
