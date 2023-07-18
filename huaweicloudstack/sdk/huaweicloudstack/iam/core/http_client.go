package core

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/converter"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/def"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/exchange"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/impl"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/request"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/response"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/sdkerr"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"net"
	"net/url"
	"reflect"
	"strings"
	"sync/atomic"
)

const (
	userAgent       = "User-Agent"
	xRequestId      = "X-Request-Id"
	contentType     = "Content-Type"
	applicationXml  = "application/xml"
	applicationBson = "application/bson"
)

type HttpClient struct {
	endpoints     []string
	endpointIndex int32
	extraHeader   map[string]string
	httpClient    *impl.DefaultHttpClient
}

func NewHttpClient(httpClient *impl.DefaultHttpClient) *HttpClient {
	return &HttpClient{httpClient: httpClient}
}

func (hc *HttpClient) WithEndpoints(endpoints []string) *HttpClient {
	hc.endpoints = endpoints
	return hc
}

func (hc *HttpClient) PreInvoke(headers map[string]string) *HttpClient {
	hc.extraHeader = headers
	return hc
}

func (hc *HttpClient) Sync(req interface{}, reqDef *def.HttpRequestDef) (*response.DefaultHttpResponse, error) {
	exg := &exchange.SdkExchange{
		ApiReference: &exchange.ApiReference{},
		Attributes:   make(map[string]interface{}),
	}
	return hc.SyncInvoke(req, reqDef, exg)
}

func (hc *HttpClient) SyncInvoke(req interface{}, reqDef *def.HttpRequestDef,
	exchange *exchange.SdkExchange) (*response.DefaultHttpResponse, error) {
	var resp *response.DefaultHttpResponse
	for {
		httpRequest, err := hc.buildRequest(req, reqDef)
		if err != nil {
			return nil, err
		}

		resp, err = hc.httpClient.SyncInvokeHttpWithExchange(httpRequest, exchange)
		if err == nil {
			break
		}

		if isNoSuchHostErr(err) && atomic.LoadInt32(&hc.endpointIndex) < int32(len(hc.endpoints)-1) {
			atomic.AddInt32(&hc.endpointIndex, 1)
		} else {
			return nil, err
		}
	}

	return hc.extractResponse(resp, reqDef)
}

func (hc *HttpClient) extractEndpoint(req interface{}, reqDef *def.HttpRequestDef, attrMaps map[string]string) (string, error) {
	var endpoint string
	for _, v := range reqDef.RequestFields {
		if v.LocationType == def.Cname {
			u, err := url.Parse(hc.endpoints[atomic.LoadInt32(&hc.endpointIndex)])
			if err != nil {
				return "", err
			}
			value, err := hc.getFieldValueByName(v.Name, attrMaps, req)
			if err != nil {
				return "", err
			}
			endpoint = fmt.Sprintf("%s://%s.%s", u.Scheme, value, u.Host)
		}
	}

	if endpoint == "" {
		endpoint = hc.endpoints[hc.endpointIndex]
	}

	return endpoint, nil
}

func (hc *HttpClient) buildRequest(req interface{}, reqDef *def.HttpRequestDef) (*request.DefaultHttpRequest, error) {
	t := reflect.TypeOf(req)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	attrMaps := hc.getFieldJsonTags(t)

	endpoint, err := hc.extractEndpoint(req, reqDef, attrMaps)
	if err != nil {
		return nil, err
	}

	builder := request.NewHttpRequestBuilder().
		WithEndpoint(endpoint).
		WithMethod(reqDef.Method).
		WithPath(reqDef.Path)

	if reqDef.ContentType != "" {
		builder.AddHeaderParam(contentType, reqDef.ContentType)
	}

	uaValue := "terraform-provider-huaweicloudstack"
	for k, v := range hc.extraHeader {
		if strings.ToLower(k) == strings.ToLower(userAgent) {
			uaValue = uaValue + ";" + v
		} else {
			builder.AddHeaderParam(k, v)
		}
	}
	builder.AddHeaderParam(userAgent, uaValue)

	builder, err = hc.fillParamsFromReq(req, t, reqDef, attrMaps, builder)
	if err != nil {
		return nil, err
	}

	var httpRequest = builder.Build()

	return httpRequest, err
}

func (hc *HttpClient) fillParamsFromReq(req interface{}, t reflect.Type, reqDef *def.HttpRequestDef,
	attrMaps map[string]string, builder *request.HttpRequestBuilder) (*request.HttpRequestBuilder, error) {

	for _, fieldDef := range reqDef.RequestFields {
		value, err := hc.getFieldValueByName(fieldDef.Name, attrMaps, req)
		if err != nil {
			return nil, err
		}

		if !value.IsValid() {
			continue
		}

		v, err := flattenEnumStruct(value)
		if err != nil {
			return nil, err
		}

		switch fieldDef.LocationType {
		case def.Header:
			builder.AddHeaderParam(fieldDef.JsonTag, fmt.Sprintf("%v", v))
		case def.Path:
			builder.AddPathParam(fieldDef.JsonTag, fmt.Sprintf("%v", v))
		case def.Query:
			builder.AddQueryParam(fieldDef.JsonTag, v)
		case def.Body:
			if body, ok := t.FieldByName("Body"); ok {
				builder.WithBody(body.Tag.Get("type"), value.Interface())
			} else {
				builder.WithBody("", value.Interface())
			}
		case def.Form:
			builder.AddFormParam(fieldDef.JsonTag, value.Interface().(def.FormData))
		}
	}

	return builder, nil
}

func (hc *HttpClient) getFieldJsonTags(t reflect.Type) map[string]string {
	attrMaps := make(map[string]string)

	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		jsonTag := t.Field(i).Tag.Get("json")
		if jsonTag != "" {
			attrMaps[t.Field(i).Name] = jsonTag
		}
	}

	return attrMaps
}

func (hc *HttpClient) getFieldValueByName(name string, jsonTag map[string]string,
	structName interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(structName)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	value := v.FieldByName(name)
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			if strings.Contains(jsonTag[name], "omitempty") {
				return reflect.ValueOf(nil), nil
			}
			return reflect.ValueOf(nil), errors.New("request field " + name + " read null value")
		}
		return value.Elem(), nil
	}

	return value, nil
}

func flattenEnumStruct(value reflect.Value) (reflect.Value, error) {
	if value.Kind() == reflect.Struct {
		v, e := jsoniter.Marshal(value.Interface())
		if e == nil {
			if strings.HasPrefix(string(v), "\"") {
				return reflect.ValueOf(strings.Trim(string(v), "\"")), nil
			} else {
				return reflect.ValueOf(string(v)), nil
			}
		}
		return reflect.ValueOf(nil), e
	}
	return value, nil
}

func (hc *HttpClient) extractResponse(resp *response.DefaultHttpResponse, reqDef *def.HttpRequestDef) (*response.DefaultHttpResponse,
	error) {
	if resp.GetStatusCode() >= 400 {
		return nil, sdkerr.NewServiceResponseError(resp.Response)
	}

	if err := hc.deserializeResponse(resp, reqDef); err != nil {

		return nil, err
	}

	return resp, nil
}

func (hc *HttpClient) deserializeResponse(resp *response.DefaultHttpResponse, reqDef *def.HttpRequestDef) error {
	t := reflect.TypeOf(reqDef.Response)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	v := reflect.ValueOf(reqDef.Response)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	addStatusCode := func() {
		field := v.FieldByName("HttpStatusCode")
		field.Set(reflect.ValueOf(resp.GetStatusCode()))
	}

	if body, ok := t.FieldByName("Body"); ok && body.Type.Name() == "ReadCloser" {
		v.FieldByName("Body").Set(reflect.ValueOf(resp.Response.Body))
		addStatusCode()
		return nil
	}

	err := hc.deserializeResponseFields(resp, reqDef)
	if err != nil {
		return err
	}

	addStatusCode()
	return nil
}

func (hc *HttpClient) deserializeResponseFields(resp *response.DefaultHttpResponse, reqDef *def.HttpRequestDef) error {
	data, err := ioutil.ReadAll(resp.Response.Body)
	if err != nil {
		if closeErr := resp.Response.Body.Close(); closeErr != nil {
			return err
		}
		return err
	}
	if err = resp.Response.Body.Close(); err != nil {
		return err
	} else {
		resp.Response.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	}

	processError := func(err error) error {
		return &sdkerr.ServiceResponseError{
			StatusCode:   resp.GetStatusCode(),
			RequestId:    resp.GetHeader(xRequestId),
			ErrorMessage: err.Error(),
		}
	}

	hasBody := false
	for _, item := range reqDef.ResponseFields {
		if item.LocationType == def.Header {
			headerErr := hc.deserializeResponseHeaders(resp, reqDef, item)
			if headerErr != nil {
				return processError(headerErr)
			}
		}

		if item.LocationType == def.Body {
			hasBody = true

			bodyErr := hc.deserializeResponseBody(reqDef, data)
			if bodyErr != nil {
				return processError(err)
			}
		}
	}

	if len(data) != 0 && !hasBody {
		if strings.Contains(resp.Response.Header.Get(contentType), applicationXml) {
			err = xml.Unmarshal(data, &reqDef.Response)
		} else if strings.Contains(resp.Response.Header.Get(contentType), applicationBson) {
			err = bson.Unmarshal(data, reqDef.Response)
		} else {
			err = jsoniter.Unmarshal(data, &reqDef.Response)
		}

		if err != nil {
			return processError(err)
		}
	}

	return nil
}

func (hc *HttpClient) deserializeResponseBody(reqDef *def.HttpRequestDef, data []byte) error {
	dataStr := string(data)

	v := reflect.ValueOf(reqDef.Response)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := reflect.TypeOf(reqDef.Response)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if body, ok := t.FieldByName("Body"); ok {
		if body.Type.Kind() == reflect.Ptr && body.Type.Elem().Kind() == reflect.String {
			v.FieldByName("Body").Set(reflect.ValueOf(&dataStr))
		} else if body.Type.Kind() == reflect.String {
			v.FieldByName("Body").Set(reflect.ValueOf(dataStr))
		} else {
			var bodyIns interface{}
			if body.Type.Kind() == reflect.Ptr {
				bodyIns = reflect.New(body.Type.Elem()).Interface()
			} else {
				bodyIns = reflect.New(body.Type).Interface()
			}
			var err error
			if reqDef.ContentType == applicationBson {
				err = bson.Unmarshal(data, bodyIns)
			} else {
				err = json.Unmarshal(data, bodyIns)
			}
			if err != nil {
				return err
			}

			if body.Type.Kind() == reflect.Ptr {
				v.FieldByName("Body").Set(reflect.ValueOf(bodyIns))
			} else {
				v.FieldByName("Body").Set(reflect.ValueOf(bodyIns).Elem())
			}
		}
	}

	return nil
}

func (hc *HttpClient) deserializeResponseHeaders(resp *response.DefaultHttpResponse, reqDef *def.HttpRequestDef,
	item *def.FieldDef) error {
	isPtr, fieldKind := hc.getFieldInfo(reqDef, item)
	v := reflect.ValueOf(reqDef.Response)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	fieldValue := v.FieldByName(item.Name)
	headerValue := resp.GetHeader(item.JsonTag)
	if headerValue == "" {
		return nil
	}

	sdkConverter := converter.StringConverterFactory(fieldKind)
	if sdkConverter == nil {
		return fmt.Errorf("failed to convert %s", item.JsonTag)
	}

	if err := sdkConverter.CovertStringToPrimitiveTypeAndSetField(fieldValue, headerValue, isPtr); err != nil {

		return err
	}

	return nil
}

func (hc *HttpClient) getFieldInfo(reqDef *def.HttpRequestDef, item *def.FieldDef) (bool, string) {

	var isPtr = false
	var fieldKind string

	t := reflect.TypeOf(reqDef.Response)
	if t.Kind() == reflect.Ptr {
		isPtr = true
		t = t.Elem()
	}

	field, _ := t.FieldByName(item.Name)
	if field.Type.Kind() == reflect.Ptr {
		fieldKind = field.Type.Elem().Kind().String()
	} else {
		fieldKind = field.Type.Kind().String()
	}

	return isPtr, fieldKind
}

func isNoSuchHostErr(err error) bool {
	if err == nil {
		return false
	}
	var errInterface interface{} = err
	if innerErr, ok := errInterface.(*url.Error); !ok {
		return false
	} else {
		errInterface = innerErr.Err
	}

	if innerErr, ok := errInterface.(*net.OpError); !ok {
		return false
	} else {
		errInterface = innerErr.Err
	}

	if innerErr, ok := errInterface.(*net.DNSError); !ok {
		return false
	} else {
		return innerErr.Err == "no such host"
	}
}
