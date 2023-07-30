// Package server implements the ricrob server.
package server

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net"
	"net/http"

	// Add profiling.
	_ "net/http/pprof"

	"github.com/go-ricrob/server/internal/exec"
)

// content is our static web server content.
//
//go:embed assets
var embeddedFS embed.FS

// Server represents a ricrob server.
type Server struct {
	svr *http.Server
}

// New creates a new server instance.
func New(logger *log.Logger, host, port string, execer *exec.Execer, solvers []string) (*Server, error) {
	rootFS, err := fs.Sub(embeddedFS, "assets")
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/board", boardHandler)
	mux.Handle("/solve", &solveHandler{execer: execer})
	mux.HandleFunc("/favicon.ico", func(http.ResponseWriter, *http.Request) {}) // Avoid "/" handler call for browser favicon request.
	mux.Handle("/", http.FileServer(http.FS(rootFS)))

	addr := net.JoinHostPort(host, port)
	return &Server{svr: &http.Server{Addr: addr, Handler: mux}}, nil
}

// ListenAndServe starts the server listening and serving content.
func (s *Server) ListenAndServe(errCh chan<- error) {
	go func() {
		if err := s.svr.ListenAndServe(); err != http.ErrServerClosed {
			errCh <- err
		}

	}()
}

// Shutdown shuts down the server gracefully.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.svr.Shutdown(ctx)
}
