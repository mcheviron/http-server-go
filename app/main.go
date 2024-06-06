package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/request"
	"github.com/codecrafters-io/http-server-starter-go/app/response"
	"github.com/codecrafters-io/http-server-starter-go/app/server"
)

func main() {
	var directory string
	args := os.Args
	for i, arg := range args {
		if arg == "--directory" && i+1 < len(args) {
			directory = args[i+1]
			break
		}
	}

	s := server.New("localhost", "4221")

	s.Get("/", func(req request.HttpRequest) response.HttpResponse {
		return response.New(response.Ok, nil, response.None)
	})

	s.Get("/echo/{str}", func(req request.HttpRequest) response.HttpResponse {
		str := req.Params["str"]
		content := &response.Content{
			Type: response.PlainText,
			Data: []byte(str),
		}
		encoding := response.None
		if strings.Contains(req.Headers["Accept-Encoding"], "gzip") {
			encoding = response.Gzip
		}
		resp := response.New(response.Ok, content, encoding)
		return resp
	})

	s.Get("/user-agent", func(req request.HttpRequest) response.HttpResponse {
		userAgent := req.Headers["User-Agent"]
		content := &response.Content{
			Type: response.PlainText,
			Data: []byte(userAgent),
		}
		resp := response.New(response.Ok, content, response.None)
		return resp
	})

	if directory != "" {
		s.Get("/files/{filename}", func(req request.HttpRequest) response.HttpResponse {
			filename := req.Params["filename"]
			filePath := filepath.Join(directory, filename)
			data, err := os.ReadFile(filePath)
			if err != nil {
				return response.New(response.NotFound, nil, response.None)
			}
			content := &response.Content{
				Type: response.OctetStream,
				Data: data,
			}
			return response.New(response.Ok, content, response.None)
		})

		s.Post("/files/{filename}", func(req request.HttpRequest) response.HttpResponse {
			filename := req.Params["filename"]
			filePath := filepath.Join(directory, filename)
			err := os.WriteFile(filePath, []byte(req.Body), 0644)
			if err != nil {
				return response.New(response.NotFound, nil, response.None)
			}
			return response.New(response.Created, nil, response.None)
		})
	}

	s.Run()
}
