package request

import "bytes"

type HttpRequest interface {
	Builder() *HttpRequestBuilder
	GetEndpoint() string
	GetMethod() string
	GetPath() string
	GetHeaderParams() map[string]string
	GetPathPrams() map[string]string
	GetQueryParams() map[string]interface{}
	GetBody() interface{}
	GetBodyToBytes() (*bytes.Buffer, error)
}
