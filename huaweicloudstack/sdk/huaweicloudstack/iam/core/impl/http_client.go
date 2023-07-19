package impl

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/request"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/response"
)

type HttpClient interface {
	sync(request *request.HttpRequest) (*response.HttpResponse, error)
}
