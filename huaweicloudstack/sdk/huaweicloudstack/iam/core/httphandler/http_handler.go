package httphandler

import (
	"net/http"
	"time"
)

type MonitorMetric struct {
	Host          string
	Path          string
	Method        string
	Raw           string
	UserAgent     string
	RequestId     string
	StatusCode    int
	ContentLength int64
	Latency       time.Duration
}

type HttpHandler struct {
	RequestHandlers  func(http.Request)
	ResponseHandlers func(http.Response)
	MonitorHandlers  func(*MonitorMetric)
}

func NewHttpHandler() *HttpHandler {
	handler := HttpHandler{}
	return &handler
}

func (handler *HttpHandler) AddRequestHandler(requestHandler func(http.Request)) *HttpHandler {
	handler.RequestHandlers = requestHandler
	return handler
}

func (handler *HttpHandler) AddResponseHandler(responseHandler func(response http.Response)) *HttpHandler {
	handler.ResponseHandlers = responseHandler
	return handler
}

func (handler *HttpHandler) AddMonitorHandler(monitorHandler func(*MonitorMetric)) *HttpHandler {
	handler.MonitorHandlers = monitorHandler
	return handler
}
