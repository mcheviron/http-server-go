package request

import (
	"bytes"
	"errors"
	"strings"
)

type HttpMethod int

const (
	GET HttpMethod = iota
	POST
	PUT
	DELETE
)

type HttpProtocol int

const (
	HTTP_1_1 HttpProtocol = iota
)

type HttpRequest struct {
	Method   HttpMethod
	Resource string
	Protocol HttpProtocol
	Headers  map[string]string
	Body     string
	Params   map[string]string
}

var (
	ErrInvalidHttpRequest      = errors.New("invalid HTTP request")
	ErrInvalidHttpMethod       = errors.New("invalid HTTP method")
	ErrUnsupportedHttpMethod   = errors.New("unsupported HTTP method")
	ErrInvalidHttpResource     = errors.New("invalid HTTP resource")
	ErrInvalidHttpProtocol     = errors.New("invalid HTTP protocol")
	ErrUnsupportedHttpProtocol = errors.New("unsupported HTTP protocol")
	ErrInvalidHttpHeader       = errors.New("invalid HTTP header")
)

func New(request []byte) (*HttpRequest, error) {
	lines := bytes.Split(request, []byte("\r\n"))
	if len(lines) < 1 {
		return nil, ErrInvalidHttpRequest
	}

	firstLineParts := bytes.Split(lines[0], []byte(" "))
	if len(firstLineParts) < 3 {
		return nil, ErrInvalidHttpRequest
	}

	methodStr := string(firstLineParts[0])
	method, err := parseHttpMethod(methodStr)
	if err != nil {
		return nil, err
	}

	resource := string(firstLineParts[1])
	if resource == "" {
		return nil, ErrInvalidHttpResource
	}

	protocolStr := string(firstLineParts[2])
	protocol, err := parseHttpProtocol(protocolStr)
	if err != nil {
		return nil, err
	}

	headers := make(map[string]string)
	for _, line := range lines[1:] {
		if len(line) == 0 {
			break
		}
		headerParts := bytes.SplitN(line, []byte(":"), 2)
		if len(headerParts) != 2 {
			return nil, ErrInvalidHttpHeader
		}
		headerName := strings.TrimSpace(string(headerParts[0]))
		headerValue := strings.TrimSpace(string(headerParts[1]))
		headers[headerName] = headerValue
	}

	body := string(bytes.Join(lines[len(headers)+2:], []byte("\r\n")))

	return &HttpRequest{
		Method:   method,
		Resource: resource,
		Protocol: protocol,
		Headers:  headers,
		Body:     body,
		Params:   nil,
	}, nil
}


func parseHttpMethod(methodStr string) (HttpMethod, error) {
	switch methodStr {
	case "GET":
		return GET, nil
	case "POST":
		return POST, nil
	case "PUT":
		return PUT, nil
	case "DELETE":
		return DELETE, nil
	default:
		return 0, ErrUnsupportedHttpMethod
	}
}

func parseHttpProtocol(protocolStr string) (HttpProtocol, error) {
	switch protocolStr {
	case "HTTP/1.1":
		return HTTP_1_1, nil
	default:
		return 0, ErrUnsupportedHttpProtocol
	}
}

