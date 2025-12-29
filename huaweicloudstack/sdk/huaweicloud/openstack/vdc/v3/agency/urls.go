package agency

import golangsdk "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const resourcePath = "rest/vdc/v3.0"

func createAgencyURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath, "vdc-agencies")
}

func getAgencyURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath, "tenant-agencies/agency-detail")
}

func listAgencyURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL(resourcePath, "vdc-agencies")
}

func getAgencyRoleURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, "vdc-agencies", id, "roles")
}

func deleteAgencyURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, "tenant-agencies", id)
}

func createAgencyProjectRoleURL(c *golangsdk.ServiceClient, agencyId, projectId, roleId string) string {
	return c.ServiceURL(resourcePath, "vdc-agencies", agencyId, "projects", projectId, "roles", roleId)
}

func createAgencyDomainRoleURL(c *golangsdk.ServiceClient, agencyId, domainId, roleId string) string {
	return c.ServiceURL(resourcePath, "vdc-agencies", agencyId, "domains", domainId, "roles", roleId)
}

func createAgencyDomainRoleInheritedURL(c *golangsdk.ServiceClient, agencyId, domainId, roleId string) string {
	return c.ServiceURL(resourcePath, "vdc-agencies", agencyId, "domains", domainId, "roles", roleId, "inherited_to_projects")
}
