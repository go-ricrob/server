package main

import (
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
	var host, port, solverList string

	//f := new(flags)
	//fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	//fmt.Fprintf(fs.Output(), "%s\n", ASCIIArt)

	// flags
	flag.StringVar(&host, fnHost, getStringEnv(envHost, "localhost"), usage(fnHost, envHost))
	flag.StringVar(&port, fnPort, getStringEnv(envPort, "50000"), usage(fnPort, envHost))
	flag.StringVar(&solverList, fnSolvers, getStringEnv(envSolvers, ""), usage(fnSolvers, envSolvers))
	flag.Parse()

	logger := log.Default()
	solvers := strings.Split(solverList, ",")

	logger.Printf("Runtime Info - GOMAXPROCS %d NumCPU %d \n", runtime.GOMAXPROCS(0), runtime.NumCPU())
	logger.Printf("Solvers %v \n", solvers)

	server := server.New(logger, host, port, solvers)
	server.ListenAndServe()
}
