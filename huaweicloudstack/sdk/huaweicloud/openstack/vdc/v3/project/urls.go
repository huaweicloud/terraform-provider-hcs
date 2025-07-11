package project

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"strings"
)

const vdcResourceBasePath = "rest/vdc"

func CreateVdcProjectURL(httpClient *golangsdk.ServiceClient) string {
	return httpClient.ServiceURL(vdcResourceBasePath, "v3.1/projects")
}

// IsValidProjectId check project valid and prevent decoded attack
func IsValidProjectId(projectId string) bool {

	if strings.Contains(projectId, "/") || strings.Contains(projectId, "..") {
		return false
	}

	return golangsdk.CheckUrlParamsValidByLoopDecode(projectId)
}

func getVdcProjectURLByProjectId(httpClient *golangsdk.ServiceClient, version string, projectId string) (string, error) {
	return httpClient.ServiceURL(vdcResourceBasePath, version, "projects", projectId), nil
}

func DeleteVdcProjectURL(httpClient *golangsdk.ServiceClient, projectId string) (string, error) {
	return getVdcProjectURLByProjectId(httpClient, "v3.0", projectId)
}

func UpdateVdcProjectURL(httpClient *golangsdk.ServiceClient, projectId string) (string, error) {
	return getVdcProjectURLByProjectId(httpClient, "v3.0", projectId)
}

func GetVdcProjectURL(httpClient *golangsdk.ServiceClient, projectId string) (string, error) {
	return getVdcProjectURLByProjectId(httpClient, "v3.1", projectId)
}
