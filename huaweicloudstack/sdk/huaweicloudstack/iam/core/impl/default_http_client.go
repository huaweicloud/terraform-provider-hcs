package impl

import (
	"bytes"
	"crypto/tls"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/config"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/exchange"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/httphandler"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/request"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/response"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type DefaultHttpClient struct {
	httpHandler  *httphandler.HttpHandler
	httpConfig   *config.HttpConfig
	transport    *http.Transport
	goHttpClient *http.Client
}

func NewDefaultHttpClient(httpConfig *config.HttpConfig) *DefaultHttpClient {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: httpConfig.IgnoreSSLVerification},
	}

	if httpConfig.DialContext != nil {
		transport.DialContext = httpConfig.DialContext
	}

	if httpConfig.HttpProxy != nil {
		proxyUrl := httpConfig.HttpProxy.GetProxyUrl()
		proxy, err := url.Parse(proxyUrl)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxy)
		}
	}

	client := &DefaultHttpClient{
		transport:  transport,
		httpConfig: httpConfig,
	}

	client.goHttpClient = &http.Client{
		Transport: client.transport,
		Timeout:   httpConfig.Timeout,
	}

	client.httpHandler = httpConfig.HttpHandler

	return client
}

func (client *DefaultHttpClient) SyncInvokeHttp(request *request.DefaultHttpRequest) (*response.DefaultHttpResponse,
	error) {
	exch := &exchange.SdkExchange{
		ApiReference: &exchange.ApiReference{},
		Attributes:   make(map[string]interface{}),
	}
	return client.SyncInvokeHttpWithExchange(request, exch)
}

func (client *DefaultHttpClient) SyncInvokeHttpWithExchange(request *request.DefaultHttpRequest,
	exch *exchange.SdkExchange) (*response.DefaultHttpResponse, error) {
	req, err := request.ConvertRequest()
	if err != nil {
		return nil, err
	}

	if lnErr := client.listenRequest(req); lnErr != nil {
		return nil, lnErr
	}

	client.recordRequestInfo(exch, req)
	resp, err := client.goHttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	client.recordResponseInfo(exch, resp)

	if lnErr := client.listenResponse(resp); lnErr != nil {
		return nil, lnErr
	}
	client.monitorHttp(exch, resp)

	return response.NewDefaultHttpResponse(resp), nil
}

func (client *DefaultHttpClient) recordRequestInfo(exch *exchange.SdkExchange, req *http.Request) {
	exch.ApiReference.Host = req.URL.Host
	exch.ApiReference.Method = req.Method
	exch.ApiReference.Path = req.URL.Path
	exch.ApiReference.Raw = req.URL.RawQuery
	exch.ApiReference.UserAgent = req.UserAgent()
	exch.ApiReference.StartedTime = time.Now()
}

func (client *DefaultHttpClient) recordResponseInfo(exch *exchange.SdkExchange, resp *http.Response) {
	exch.ApiReference.RequestId = resp.Header.Get("X-Request-Id")
	exch.ApiReference.StatusCode = resp.StatusCode
	exch.ApiReference.ContentLength = resp.ContentLength
	exch.ApiReference.DurationMs = time.Since(exch.ApiReference.StartedTime)
}

func (client *DefaultHttpClient) listenRequest(req *http.Request) error {
	if client.httpHandler != nil && client.httpHandler.RequestHandlers != nil && req != nil {
		reqClone := req.Clone(req.Context())

		if req.Body != nil {
			bodyBytes, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return err
			}

			req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			reqClone.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			defer reqClone.Body.Close()
		}

		client.httpHandler.RequestHandlers(*reqClone)
	}

	return nil
}

func (client *DefaultHttpClient) listenResponse(resp *http.Response) error {
	if client.httpHandler != nil && client.httpHandler.ResponseHandlers != nil && resp != nil {
		respClone := http.Response{
			Status:           resp.Status,
			StatusCode:       resp.StatusCode,
			Proto:            resp.Proto,
			ProtoMajor:       resp.ProtoMajor,
			ProtoMinor:       resp.ProtoMinor,
			Header:           resp.Header,
			ContentLength:    resp.ContentLength,
			TransferEncoding: resp.TransferEncoding,
			Close:            resp.Close,
			Uncompressed:     resp.Uncompressed,
			Trailer:          resp.Trailer,
		}

		if resp.Body != nil {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			respClone.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			defer respClone.Body.Close()
		}

		client.httpHandler.ResponseHandlers(respClone)
	}

	return nil
}

func (client *DefaultHttpClient) monitorHttp(exch *exchange.SdkExchange, resp *http.Response) {
	if client.httpHandler != nil && client.httpHandler.MonitorHandlers != nil {
		metric := &httphandler.MonitorMetric{
			Host:          exch.ApiReference.Host,
			Method:        exch.ApiReference.Method,
			Path:          exch.ApiReference.Path,
			Raw:           exch.ApiReference.Raw,
			UserAgent:     exch.ApiReference.UserAgent,
			Latency:       exch.ApiReference.DurationMs,
			RequestId:     exch.ApiReference.RequestId,
			StatusCode:    exch.ApiReference.StatusCode,
			ContentLength: exch.ApiReference.ContentLength,
		}

		client.httpHandler.MonitorHandlers(metric)
	}
}
