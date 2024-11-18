package snapshots

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

func actionURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL("servers", id, "action")
}

func getURL(client *golangsdk.ServiceClient, serverId string) string {
	return client.ServiceURL("images") + "?__snapshot_from_instance=" + serverId
}

func jobURL(client *golangsdk.ServiceClient, jobId string) string {
	return client.ServiceURL("jobs", jobId)
}

func deleteURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("cloudimages")
}
