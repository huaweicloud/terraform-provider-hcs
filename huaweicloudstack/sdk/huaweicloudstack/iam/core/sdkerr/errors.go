package sdkerr

import (
	"bytes"
	"fmt"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/utils"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

const (
	xRequestId                  = "X-Request-Id"
	code                        = "code"
	message                     = "message"
	errorCode                   = "error_code"
	errorMsg                    = "error_msg"
	encodedAuthorizationMessage = "encoded_authorization_message"
)

type CredentialsTypeError struct {
	ErrorMessage string
}

func NewCredentialsTypeError(msg string) *CredentialsTypeError {
	c := &CredentialsTypeError{
		ErrorMessage: msg,
	}
	return c
}

func (c *CredentialsTypeError) Error() string {
	return fmt.Sprintf("{\"ErrorMessage\": \"%s\"}", c.ErrorMessage)
}

type ConnectionError struct {
	ErrorMessage string
}

func NewConnectionError(msg string) *ConnectionError {
	c := &ConnectionError{
		ErrorMessage: msg,
	}
	return c
}

func (c *ConnectionError) Error() string {
	return fmt.Sprintf("{\"ErrorMessage\": \"%s\"}", c.ErrorMessage)
}

type RequestTimeoutError struct {
	ErrorMessage string
}

func NewRequestTimeoutError(msg string) *RequestTimeoutError {
	rt := &RequestTimeoutError{
		ErrorMessage: msg,
	}
	return rt
}

func (rt *RequestTimeoutError) Error() string {
	return fmt.Sprintf("{\"ErrorMessage\": \"%s\"}", rt.ErrorMessage)
}

type errMap map[string]interface{}

func (m errMap) getStringValue(key string) string {
	var result string

	value, isExist := m[key]
	if !isExist {
		return result
	}

	if strVal, ok := value.(string); ok {
		result = strVal
	}

	return result
}

type ServiceResponseError struct {
	StatusCode                  int    `json:"status_code"`
	RequestId                   string `json:"request_id"`
	ErrorCode                   string `json:"error_code"`
	ErrorMessage                string `json:"error_message"`
	EncodedAuthorizationMessage string `json:"encoded_authorization_message"`
}

func NewServiceResponseError(resp *http.Response) *ServiceResponseError {
	sr := &ServiceResponseError{
		StatusCode: resp.StatusCode,
		RequestId:  resp.Header.Get(xRequestId),
	}

	data, err := ioutil.ReadAll(resp.Body)
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr == nil && err == nil {
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		}
	}()

	if err == nil {
		dataBuf := errMap{}
		err := jsoniter.Unmarshal(data, &dataBuf)
		if err != nil {
			sr.ErrorMessage = string(data)
		} else {
			processServiceResponseError(dataBuf, sr)
			if sr.ErrorMessage == "" {
				sr.ErrorMessage = string(data)
			}
		}
	}

	return sr
}

func processServiceResponseError(m errMap, sr *ServiceResponseError) {
	if value := m.getStringValue(encodedAuthorizationMessage); value != "" {
		sr.EncodedAuthorizationMessage = value
	}

	_code := m.getStringValue(errorCode)
	msg := m.getStringValue(errorMsg)
	if _code != "" && msg != "" {
		sr.ErrorCode = _code
		sr.ErrorMessage = msg
		return
	}

	_code = m.getStringValue(code)
	msg = m.getStringValue(message)
	if _code != "" && msg != "" {
		sr.ErrorCode = _code
		sr.ErrorMessage = msg
		return
	}

	for _, v := range m {
		if val, ok := v.(map[string]interface{}); ok {
			processServiceResponseError(val, sr)
		}
	}
}

func (sr ServiceResponseError) Error() string {
	data, err := utils.Marshal(sr)
	if err != nil {
		return fmt.Sprintf("{\"ErrorMessage\": \"%s\",\"ErrorCode\": \"%s\",\"EncodedAuthorizationMessage\": \"%s\"}",
			sr.ErrorMessage, sr.ErrorCode, sr.EncodedAuthorizationMessage)
	}
	return fmt.Sprintf(string(data))
}

type BadRequestError struct {
	ErrorMessage string
}

func NewBadRequestError(msg string) *BadRequestError {
	c := &BadRequestError{
		ErrorMessage: msg,
	}
	return c
}

func (br BadRequestError) Error() string {
	return fmt.Sprintf("{\"ErrorMessage\": \"%s\"}", br.ErrorMessage)
}
