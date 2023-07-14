package flavors

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("baremetalservers", "flavors")
}
