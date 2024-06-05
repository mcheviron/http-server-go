package server

import (
	"log/slog"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/request"
	"github.com/codecrafters-io/http-server-starter-go/app/response"
)

type Route struct {
	method  request.HttpMethod
	path    string
	params  []string
	handler func(request.HttpRequest) response.HttpResponse
}

type routeKey struct {
	method request.HttpMethod
	path   string
}

type Router struct {
	routes map[routeKey]Route
}

func (r *Router) addRoute(method request.HttpMethod, path string, handler func(request.HttpRequest) response.HttpResponse) {
	params := extractParams(path)
	key := routeKey{method, path}
	r.routes[key] = Route{method, path, params, handler}
}

func (r *Router) Get(path string, handler func(request.HttpRequest) response.HttpResponse) {
	r.addRoute(request.GET, path, handler)
}

func (r *Router) Post(path string, handler func(request.HttpRequest) response.HttpResponse) {
	r.addRoute(request.POST, path, handler)
}

func (r *Router) Put(path string, handler func(request.HttpRequest) response.HttpResponse) {
	r.addRoute(request.PUT, path, handler)
}

func (r *Router) Delete(path string, handler func(request.HttpRequest) response.HttpResponse) {
	r.addRoute(request.DELETE, path, handler)
}

func (r *Router) handleRequest(httpRequest request.HttpRequest) response.HttpResponse {
	for key, route := range r.routes {
		if key.method == httpRequest.Method {
			if params := matchRoute(route.path, httpRequest.Resource); params != nil {
				httpRequest.Params = params
				return route.handler(httpRequest)
			}
		}
	}
	return *response.New(response.NotFound, nil, nil)
}

func extractParams(path string) []string {
	var params []string
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			params = append(params, part[1:len(part)-1])
		}
	}
	return params
}

func matchRoute(routePath, requestPath string) map[string]string {
	routeParts := strings.Split(routePath, "/")
	requestParts := strings.Split(requestPath, "/")

	if len(routeParts) != len(requestParts) {
		return nil
	}

	params := make(map[string]string)
	for i := range routeParts {
		if strings.HasPrefix(routeParts[i], "{") && strings.HasSuffix(routeParts[i], "}") {
			paramName := routeParts[i][1 : len(routeParts[i])-1]
			params[paramName] = requestParts[i]
		} else if routeParts[i] != requestParts[i] {
			return nil
		}
	}
	return params
}

type Server struct {
	Router
	host string
	port string
}

func New(host string, port string) *Server {
	return &Server{
		Router: Router{
			routes: make(map[routeKey]Route),
		},
		host: host,
		port: port,
	}
}

func (s *Server) Run() {
	slog.Info("Logs from your program will appear here!")

	address := net.JoinHostPort(s.host, s.port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		slog.Error("Failed to bind to address", slog.String("error", err.Error()))
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			slog.Error("Error accepting connection", slog.String("error", err.Error()))
			os.Exit(1)
		}

		go s.handleConnection(conn)
	}
}
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		slog.Error("Error reading request", slog.String("error", err.Error()))
		return
	}

	httpRequest, err := request.New(buf)
	if err != nil {
		slog.Error("Error parsing request", slog.String("error", err.Error()))
		httpResponse := response.New(response.NotFound, nil, nil)
		_, err = conn.Write(httpResponse.Bytes())
		if err != nil {
			slog.Error("Error writing response", slog.String("error", err.Error()))
		}
		return
	}

	httpResponse := s.handleRequest(*httpRequest)
	_, err = conn.Write(httpResponse.Bytes())
	if err != nil {
		slog.Error("Error writing response", slog.String("error", err.Error()))
	}
}
