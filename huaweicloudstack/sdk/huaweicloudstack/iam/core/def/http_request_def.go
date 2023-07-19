package def

type HttpRequestDef struct {
	Method         string
	Path           string
	ContentType    string
	RequestFields  []*FieldDef
	ResponseFields []*FieldDef
	Response       interface{}
}

type HttpRequestDefBuilder struct {
	httpRequestDef *HttpRequestDef
}

func NewHttpRequestDefBuilder() *HttpRequestDefBuilder {
	httpRequestDef := &HttpRequestDef{
		RequestFields:  []*FieldDef{},
		ResponseFields: []*FieldDef{},
	}
	HttpRequestDefBuilder := &HttpRequestDefBuilder{
		httpRequestDef: httpRequestDef,
	}
	return HttpRequestDefBuilder
}

func (builder *HttpRequestDefBuilder) WithPath(path string) *HttpRequestDefBuilder {
	builder.httpRequestDef.Path = path
	return builder
}

func (builder *HttpRequestDefBuilder) WithMethod(method string) *HttpRequestDefBuilder {
	builder.httpRequestDef.Method = method
	return builder
}

func (builder *HttpRequestDefBuilder) WithContentType(contentType string) *HttpRequestDefBuilder {
	builder.httpRequestDef.ContentType = contentType
	return builder
}

func (builder *HttpRequestDefBuilder) WithResponse(response interface{}) *HttpRequestDefBuilder {
	builder.httpRequestDef.Response = response
	return builder
}

func (builder *HttpRequestDefBuilder) WithRequestField(field *FieldDef) *HttpRequestDefBuilder {
	builder.httpRequestDef.RequestFields = append(builder.httpRequestDef.RequestFields, field)
	return builder
}

func (builder *HttpRequestDefBuilder) WithResponseField(field *FieldDef) *HttpRequestDefBuilder {
	builder.httpRequestDef.ResponseFields = append(builder.httpRequestDef.ResponseFields, field)
	return builder
}

func (builder *HttpRequestDefBuilder) Build() *HttpRequestDef {
	return builder.httpRequestDef
}
