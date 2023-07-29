package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

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

	server := server.New(logger, host, port, solvers)
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}
