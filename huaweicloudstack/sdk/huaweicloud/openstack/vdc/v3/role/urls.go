package role

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"strings"
)

const vdcResourceBasePath = "rest/vdc/v3.0"

func ListVdcRoleURL(httpClient *golangsdk.ServiceClient) string {
	return httpClient.ServiceURL(vdcResourceBasePath, "OS-ROLE/roles/third-party/roles")
}

func CreateVdcRoleURL(httpClient *golangsdk.ServiceClient) string {
	return httpClient.ServiceURL(vdcResourceBasePath, "OS-ROLE/roles")
}

// IsValidRoleId check roleId valid and prevent decoded attack
func IsValidRoleId(roleId string) bool {

	if strings.Contains(roleId, "/") || strings.Contains(roleId, "..") {
		return false
	}

	return golangsdk.CheckUrlParamsValidByLoopDecode(roleId)
}

func getVdcRoleURLByRoleId(httpClient *golangsdk.ServiceClient, roleId string) (string, error) {
	return httpClient.ServiceURL(vdcResourceBasePath, "OS-ROLE/roles", roleId), nil
}

func DeleteVdcRoleURL(httpClient *golangsdk.ServiceClient, roleId string) (string, error) {
	return getVdcRoleURLByRoleId(httpClient, roleId)
}

func UpdateVdcRoleURL(httpClient *golangsdk.ServiceClient, roleId string) (string, error) {
	return getVdcRoleURLByRoleId(httpClient, roleId)
}

func GetVdcRoleURL(httpClient *golangsdk.ServiceClient, roleId string) (string, error) {
	return getVdcRoleURLByRoleId(httpClient, roleId)
}
