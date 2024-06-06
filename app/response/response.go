package response

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

type ResponseType int

const (
	Ok ResponseType = iota
	Created
	NotFound
	InternalServerError
)

type ContentType int

const (
	PlainText ContentType = iota
	OctetStream
)

type Encoding int

const (
	None Encoding = iota
	Gzip
)

type Content struct {
	Type ContentType
	Data []byte
}

type HttpResponse struct {
	Type     ResponseType
	Content  *Content
	Encoding Encoding
}

func New(responseType ResponseType, content *Content, encoding Encoding) HttpResponse {
	return HttpResponse{
		Type:     responseType,
		Content:  content,
		Encoding: encoding,
	}
}

func (r *HttpResponse) Bytes() []byte {
	var buf bytes.Buffer

	switch r.Type {
	case Ok:
		buf.WriteString("HTTP/1.1 200 OK\r\n")
		if r.Content != nil {
			switch r.Content.Type {
			case PlainText:
				buf.WriteString("Content-Type: text/plain\r\n")
				r.writeBody(&buf)
			case OctetStream:
				buf.WriteString("Content-Type: application/octet-stream\r\n")
				r.writeBody(&buf)
			}
		}
	case Created:
		buf.WriteString("HTTP/1.1 201 Created\r\n")
	case NotFound:
		buf.WriteString("HTTP/1.1 404 Not Found\r\n")
	default:
		buf.WriteString("HTTP/1.1 500 Internal Server Error\r\n")
	}

	buf.WriteString("\r\n")

	return buf.Bytes()
}

func (r *HttpResponse) writeBody(buf *bytes.Buffer) {
	var body []byte
	if r.Encoding == Gzip {
		buf.WriteString("Content-Encoding: gzip\r\n")
		var compressedData bytes.Buffer
		gzipWriter := gzip.NewWriter(&compressedData)
		_, _ = gzipWriter.Write(r.Content.Data)
		_ = gzipWriter.Close()
		body = compressedData.Bytes()
	} else {
		body = r.Content.Data
	}
	buf.WriteString("Content-Length: " + fmt.Sprint(len(body)) + "\r\n")
	buf.WriteString("\r\n")
	buf.Write(body)
}
