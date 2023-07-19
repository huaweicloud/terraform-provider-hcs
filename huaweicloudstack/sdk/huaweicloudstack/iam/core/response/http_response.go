package response

type HttpResponse interface {
	GetStatusCode() int
	GetHeaders() map[string]string
	GetBody() string
	GetBodyJson() interface{}
	GetHeader(key string) string
}
