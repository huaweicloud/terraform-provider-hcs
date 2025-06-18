package role

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

const vdcResourceBasePath = "rest/vdc/v3.0"

func ListVdcRoleURL(httpClient *golangsdk.ServiceClient) string {
	return httpClient.ServiceURL(vdcResourceBasePath, "OS-ROLE/roles/third-party/roles")
}

func CreateVdcRoleURL(httpClient *golangsdk.ServiceClient) string {
	return httpClient.ServiceURL(vdcResourceBasePath, "OS-ROLE/roles")
}

func getVdcRoleURLByRoleId(httpClient *golangsdk.ServiceClient, roleId string) string {
	return httpClient.ServiceURL(vdcResourceBasePath, "OS-ROLE/roles", roleId)
}

func DeleteVdcRoleURL(httpClient *golangsdk.ServiceClient, roleId string) string {
	return getVdcRoleURLByRoleId(httpClient, roleId)
}

func UpdateVdcRoleURL(httpClient *golangsdk.ServiceClient, roleId string) string {
	return getVdcRoleURLByRoleId(httpClient, roleId)
}

func GetVdcRoleURL(httpClient *golangsdk.ServiceClient, roleId string) string {
	return getVdcRoleURLByRoleId(httpClient, roleId)
}
