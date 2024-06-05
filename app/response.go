package main

import (
	"bytes"
)

type HttpResponse struct {
	Type     ResponseType
	Content  *Content
	Encoding *Encoding
}

func (r *HttpResponse) Bytes() []byte {
	var buf bytes.Buffer

	switch r.Type {
	case Ok:
		buf.WriteString("HTTP/1.1 200 OK\r\n")
	case Created:
		buf.WriteString("HTTP/1.1 201 Created\r\n")
	case NotFound:
		buf.WriteString("HTTP/1.1 404 Not Found\r\n")
	}

	buf.WriteString("\r\n")

	return buf.Bytes()
}

type ResponseType int

const (
	Ok ResponseType = iota
	Created
	NotFound
)

type Content struct {
	Type ContentType
	Data []byte
}

type ContentType int

const (
	PlainText ContentType = iota
	OctetStream
)

type Encoding int

const (
	Gzip Encoding = iota
)
