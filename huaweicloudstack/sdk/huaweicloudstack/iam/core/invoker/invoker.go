package invoker

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/def"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/exchange"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/invoker/retry"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/response"
	"time"
)

type RetryChecker func(interface{}, error) bool

type BaseInvoker struct {
	Exchange *exchange.SdkExchange

	client  *core.HttpClient
	request interface{}
	meta    *def.HttpRequestDef
	headers map[string]string

	retryTimes      int
	retryChecker    RetryChecker
	backoffStrategy retry.Strategy
}

func NewBaseInvoker(client *core.HttpClient, request interface{}, meta *def.HttpRequestDef) *BaseInvoker {
	exch := &exchange.SdkExchange{
		ApiReference: &exchange.ApiReference{
			Method: meta.Method,
			Path:   meta.Path,
		},
		Attributes: make(map[string]interface{}),
	}

	return &BaseInvoker{
		Exchange: exch,
		client:   client,
		request:  request,
		meta:     meta,
		headers:  make(map[string]string),
	}
}

func (b *BaseInvoker) AddHeader(headers map[string]string) *BaseInvoker {
	b.headers = headers
	return b
}

func (b *BaseInvoker) WithRetry(retryTimes int, checker RetryChecker, backoffStrategy retry.Strategy) *BaseInvoker {
	b.retryTimes = retryTimes
	b.retryChecker = checker
	b.backoffStrategy = backoffStrategy
	return b
}

func (b *BaseInvoker) Invoke() (*response.DefaultHttpResponse, error) {
	if b.retryTimes != 0 && b.retryChecker != nil {
		var execTimes int
		var resp *response.DefaultHttpResponse
		var err error
		for {
			if execTimes == b.retryTimes {
				break
			}
			resp, err = b.client.PreInvoke(b.headers).SyncInvoke(b.request, b.meta, b.Exchange)
			execTimes += 1

			if b.retryChecker(resp, err) {
				time.Sleep(time.Duration(b.backoffStrategy.ComputeDelayBeforeNextRetry()))
			} else {
				break
			}
		}
		return resp, err
	} else {
		return b.client.PreInvoke(b.headers).SyncInvoke(b.request, b.meta, b.Exchange)
	}
}
