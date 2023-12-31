package maintainwindows

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

// Get maintain windows
func Get(client *golangsdk.ServiceClient) ([]MaintainWindow, error) {
	var rst golangsdk.Result
	_, err := client.Get(getURL(client), &rst.Body, nil)
	if err == nil {
		var r GetResponse
		err = rst.ExtractInto(&r)
		return r.MaintainWindows, err
	}
	return nil, err
}
