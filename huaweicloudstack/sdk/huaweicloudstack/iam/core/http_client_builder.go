package core

import (
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/impl"
	"strings"
)

type HttpClientBuilder struct {
	derivedAuthServiceName string
	endpoints              []string
	httpConfig             *config.HttpConfig
}

func NewHcsHttpClientBuilder() *HttpClientBuilder {
	hcHttpClientBuilder := &HttpClientBuilder{}
	return hcHttpClientBuilder
}

func (builder *HttpClientBuilder) WithDerivedAuthServiceName(derivedAuthServiceName string) *HttpClientBuilder {
	builder.derivedAuthServiceName = derivedAuthServiceName
	return builder
}

func (builder *HttpClientBuilder) WithEndpoints(endpoints []string) *HttpClientBuilder {
	builder.endpoints = endpoints
	return builder
}

func (builder *HttpClientBuilder) WithHttpConfig(httpConfig *config.HttpConfig) *HttpClientBuilder {
	builder.httpConfig = httpConfig
	return builder
}

func (builder *HttpClientBuilder) Build() *HttpClient {
	if builder.httpConfig == nil {
		builder.httpConfig = config.DefaultHttpConfig()
	}

	defaultHttpClient := impl.NewDefaultHttpClient(builder.httpConfig)

	for index, endpoint := range builder.endpoints {
		if !strings.HasPrefix(endpoint, "http") {
			builder.endpoints[index] = "https://" + endpoint
		}
	}

	hcHttpClient := NewHttpClient(defaultHttpClient).WithEndpoints(builder.endpoints)
	return hcHttpClient
}
