package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"

	"github.com/go-ricrob/server/internal/exec"
	"github.com/go-ricrob/server/internal/server"

	// Add profiling.
	_ "net/http/pprof"
)

//go:embed asciiart.txt
var asciiart string

func getStringEnv(key, defValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defValue
	}
	return value
}

const (
	fnHost    = "host"
	fnPort    = "port"
	fnSolvers = "solvers"
)

const (
	envHost    = "HOST"
	envPort    = "PORT"
	envSolvers = "SOLVERS"
)

func usage(name, envName string) string {
	return fmt.Sprintf("%s (environment variable: %s)", name, envName)
}

func main() {
	logger := log.Default()
	var host, port, solverList string

	fmt.Println(asciiart)

	// flags
	flag.StringVar(&host, fnHost, getStringEnv(envHost, "localhost"), usage(fnHost, envHost))
	flag.StringVar(&port, fnPort, getStringEnv(envPort, "50000"), usage(fnPort, envHost))
	flag.StringVar(&solverList, fnSolvers, getStringEnv(envSolvers, ""), usage(fnSolvers, envSolvers))
	flag.Parse()

	solvers := strings.Split(solverList, ",")

	logger.Printf("Runtime Info - GOMAXPROCS %d NumCPU %d \n", runtime.GOMAXPROCS(0), runtime.NumCPU())
	logger.Printf("Solvers %v \n", solvers)

	execer := exec.New(solvers, logger)
	server, err := server.New(logger, host, port, execer, solvers)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Printf("please open your web browser at http://%s:%s \n", host, port)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	errCh := make(chan error)
	defer close(errCh)

	server.ListenAndServe(errCh)
	select {
	case <-sigint:
		logger.Println("shutting down...")
	case err := <-errCh:
		logger.Fatal(err)
	}
	if err := server.Shutdown(context.Background()); err != nil {
		logger.Fatal(err)
	}
}
