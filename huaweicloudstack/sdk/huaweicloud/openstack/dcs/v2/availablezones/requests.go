package availablezones

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
)

func List(client *golangsdk.ServiceClient) (*GetResponse, error) {
	var rst golangsdk.Result
	_, rst.Err = client.Get(getURL(client), &rst.Body, nil)
	if rst.Err == nil {
		var s GetResponse
		err := rst.ExtractInto(&s)
		return &s, err
	}
	return nil, rst.Err
}
