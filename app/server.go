package main

import (
	"log/slog"
	"net"
	"os"
)

func main() {
	slog.Info("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		slog.Error("Failed to bind to port 4221", slog.String("error", err.Error()))
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			slog.Error("Error accepting connection", slog.String("error", err.Error()))
			os.Exit(1)
		}

		response := &HttpResponse{Type: Ok}
		_, err = conn.Write(response.Bytes())
		if err != nil {
			slog.Error("Error writing response", slog.String("error", err.Error()))
		}
		conn.Close()
	}
}
