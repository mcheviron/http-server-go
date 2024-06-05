package main

import (
	"github.com/codecrafters-io/http-server-starter-go/app/request"
	"github.com/codecrafters-io/http-server-starter-go/app/response"
	"github.com/codecrafters-io/http-server-starter-go/app/server"
)

func main() {
	s := server.New("localhost", "4221")

	s.Get("/", func(req request.HttpRequest) response.HttpResponse {
		return *response.New(response.Ok, nil, nil)
	})

	s.Get("/echo/{str}", func(req request.HttpRequest) response.HttpResponse {
		str := req.Params["str"]
		content := &response.Content{
			Type: response.PlainText,
			Data: []byte(str),
		}
		resp := response.New(response.Ok, content, nil)
		return *resp
	})

	s.Get("/user-agent", func(req request.HttpRequest) response.HttpResponse {
 		userAgent := req.Headers["User-Agent"]
 		content := &response.Content{
 			Type: response.PlainText,
 			Data: []byte(userAgent),
 		}
 		resp := response.New(response.Ok, content, nil)
 		return *resp
 	})
 
 

	s.Run()
}
