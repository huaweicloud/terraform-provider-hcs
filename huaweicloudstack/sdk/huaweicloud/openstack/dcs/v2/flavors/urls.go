package flavors

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

// listURL will build the get url of List function
func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(client.ProjectID, "flavors")
}
