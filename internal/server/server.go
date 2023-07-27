// Package server implements the ricrob server.
package server

import (
	"context"
	"embed"
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
//go:embed index.html
//go:embed assets
var content embed.FS

// const (
// 	docRoot = "www"
// )

// rootFS
type rootFS struct {
	fs http.FileSystem
}

func (fs rootFS) Open(name string) (http.File, error) {
	// TODO try to 're-route' to docRoot
	log.Println(name)
	return fs.fs.Open(name)
}

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
func (s *Server) ListenAndServe() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	execCmd := newExecCmd(s.solvers, s.logger)

	http.HandleFunc("/board", boardHandler)
	http.Handle("/solve", &solveHandler{execCmd: execCmd})
	http.HandleFunc("/favicon.ico", func(http.ResponseWriter, *http.Request) {}) // Avoid "/" handler call for browser favicon request.
	http.Handle("/", http.FileServer(rootFS{http.FS(content)}))

	addr := net.JoinHostPort(s.host, s.port)
	//svr := http.Server{Addr: addr, Handler: mux}
	svr := http.Server{Addr: addr}
	log.Printf("listening on %s ...\n", addr)

	go func() {
		if err := svr.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-sigint
	// shutdown server
	log.Println("shutting down...")

	if err := svr.Shutdown(context.Background()); err != nil {
		log.Fatalf("HTTP server Shutdown: %v", err)
	}
}
