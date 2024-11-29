package clone

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloud/openstack/ecs/v1/cloudservers"
)

type GetResult struct {
	golangsdk.Result
}

func (r GetResult) Extract() (*cloudservers.CloudServer, error) {
	var s struct {
		Server *cloudservers.CloudServer `json:"server"`
	}
	err := r.ExtractInto(&s)
	return s.Server, err
}
