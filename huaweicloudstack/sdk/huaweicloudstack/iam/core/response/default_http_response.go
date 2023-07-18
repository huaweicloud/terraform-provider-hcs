package response

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type DefaultHttpResponse struct {
	Response *http.Response
}

func NewDefaultHttpResponse(response *http.Response) *DefaultHttpResponse {
	return &DefaultHttpResponse{Response: response}
}

func (r *DefaultHttpResponse) GetStatusCode() int {
	return r.Response.StatusCode
}

func (r *DefaultHttpResponse) GetHeaders() map[string]string {
	headerParams := map[string]string{}
	for key, values := range r.Response.Header {
		if values == nil || len(values) <= 0 {
			continue
		}
		headerParams[key] = values[0]
	}
	return headerParams
}

func (r *DefaultHttpResponse) GetBody() string {
	body, err := ioutil.ReadAll(r.Response.Body)
	if err != nil {
		return ""
	}
	if err := r.Response.Body.Close(); err == nil {
		r.Response.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}
	return string(body)
}

func (r *DefaultHttpResponse) GetHeader(key string) string {
	header := r.Response.Header
	return header.Get(key)
}
