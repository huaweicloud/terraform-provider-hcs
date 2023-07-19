package def

import (
	"fmt"
	"github.com/huaweicloud/terraform-provider-hcs/huaweicloudstack/sdk/huaweicloudstack/iam/core/converter"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"strings"
)

var quoteEscape = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscape.Replace(s)
}

type FilePart struct {
	Headers textproto.MIMEHeader
	Content *os.File
}

func NewFilePart(content *os.File) *FilePart {
	return &FilePart{
		Content: content,
	}
}

func NewFilePartWithContentType(content *os.File, contentType string) *FilePart {
	var headers = make(textproto.MIMEHeader)
	headers.Set("Content-Type", contentType)

	return &FilePart{
		Headers: headers,
		Content: content,
	}
}

func (f FilePart) Write(w *multipart.Writer, name string) error {
	var h textproto.MIMEHeader
	if f.Headers != nil {
		h = f.Headers
	} else {
		h = make(textproto.MIMEHeader)
	}

	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(name), escapeQuotes(f.Content.Name())))

	if f.Headers.Get("Content-Type") == "" {
		h.Set("Content-Type", "application/octet-stream")
	}

	writer, err := w.CreatePart(h)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, f.Content)
	return err
}

type MultiPart struct {
	Content interface{}
}

func NewMultiPart(content interface{}) *MultiPart {
	return &MultiPart{
		Content: content,
	}
}

func (m MultiPart) Write(w *multipart.Writer, name string) error {
	err := w.WriteField(name, converter.ConvertInterfaceToString(m.Content))
	return err
}

type FormData interface {
	Write(*multipart.Writer, string) error
}
