package group_role_assignment

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"net/url"
)

const rootPath = "rest/vdc/v3.2"

// User Group - Permission Management
func GroupRoleAssignmentURL(c *golangsdk.ServiceClient, groupId string) string {
	return c.ServiceURL(rootPath, "groups", url.PathEscape(groupId), "roles")
}

func GroupRoleAssignmentEpsURL(c *golangsdk.ServiceClient, epId string, groupId string, roleId string) string {
	return c.ServiceURL("v1.0/enterprise-projects", url.PathEscape(epId), "groups", url.PathEscape(groupId), "roles", url.PathEscape(roleId))
}
