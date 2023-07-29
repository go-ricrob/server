// Package server implements the ricrob server.
package server

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	// Add profiling.
	_ "net/http/pprof"
)

// content is our static web server content.
//
//go:embed assets
var embeddedFS embed.FS

// Server represents a ricrob server.
type Server struct {
	logger     *log.Logger
	host, port string
	solvers    []string
}

// New creates a new server instance.
func New(logger *log.Logger, host, port string, solvers []string) *Server {
	return &Server{
		logger:  logger,
		host:    host,
		port:    port,
		solvers: solvers,
	}
}

// ListenAndServe starts the server listening and serving content.
func (s *Server) ListenAndServe() error {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	//mime.AddExtensionType(".wasm", "application/wasm")
	//mime.AddExtensionType("js", "text/javascript")
	rootFS, err := fs.Sub(embeddedFS, "assets")
	if err != nil {
		return err
	}

	execCmd := newExecCmd(s.solvers, s.logger)

	http.HandleFunc("/board", boardHandler)
	http.Handle("/solve", &solveHandler{execCmd: execCmd})
	http.HandleFunc("/favicon.ico", func(http.ResponseWriter, *http.Request) {}) // Avoid "/" handler call for browser favicon request.
	http.Handle("/", http.FileServer(http.FS(rootFS)))

	addr := net.JoinHostPort(s.host, s.port)
	//svr := http.Server{Addr: addr, Handler: mux}
	svr := http.Server{Addr: addr}
	s.logger.Printf("listening on %s ...\n", addr)

	go func() {
		if err := svr.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Fatal(err)
		}
	}()

	<-sigint
	// shutdown server
	s.logger.Println("shutting down...")

	return svr.Shutdown(context.Background())
}
