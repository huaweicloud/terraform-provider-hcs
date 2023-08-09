package baremetalservers

import "github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"

func createURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL("baremetalservers")
}

func getURL(sc *golangsdk.ServiceClient, serverID string) string {
	return sc.ServiceURL("cloudservers", serverID)
}

func putURL(sc *golangsdk.ServiceClient, serverID string) string {
	return sc.ServiceURL("baremetalservers", serverID)
}

func deleteURL(sc *golangsdk.ServiceClient) string {
	return sc.ServiceURL("baremetalservers/delete")
}

func jobURL(sc *golangsdk.ServiceClient, jobId string) string {
	return sc.ServiceURL("jobs", jobId)
}
