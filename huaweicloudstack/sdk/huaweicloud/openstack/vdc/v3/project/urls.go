package project

import (
	"fmt"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"net/url"
	"slices"
	"strings"
)

const vdcResourceBasePath = "rest/vdc"

func CreateVdcProjectURL(httpClient *golangsdk.ServiceClient) string {
	return httpClient.ServiceURL(vdcResourceBasePath, "v3.1/projects")
}

func checkUrlParamsValidByLoopDecode(pathParams string) bool {
	decoded := pathParams
	maxDecodeTimes := 3
	var i int

	for i = 0; i < maxDecodeTimes; i++ {
		unescape, err := url.PathUnescape(decoded)
		if err != nil || unescape == decoded {
			decoded = unescape
			break
		}

		decoded = unescape

		if strings.Contains(decoded, "/") || strings.Contains(decoded, "..") {
			return false
		}
	}

	return true
}

// IsValidProjectId check project valid and prevent decoded attack
func IsValidProjectId(projectId string) bool {

	if strings.Contains(projectId, "/") || strings.Contains(projectId, "..") {
		return false
	}

	return checkUrlParamsValidByLoopDecode(projectId)
}

func IsValidVersion(version string) bool {
	versions := []string{"v3.0", "v3.1"}
	return slices.Contains(versions, version)
}

func getVdcProjectURLByProjectId(httpClient *golangsdk.ServiceClient, version string, projectId string) (string, error) {
	if IsValidProjectId(projectId) && IsValidVersion(version) {
		return httpClient.ServiceURL(vdcResourceBasePath, version, "projects", url.PathEscape(projectId)), nil
	} else {
		return "", fmt.Errorf("invalid project_id or version")
	}
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
