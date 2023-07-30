// Package exec runs solvers.
package exec

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"sync"
)

// Result represents a solver result.
type Result struct {
	Response []byte
	Err      error
}

// Execer represents the solvers to be run.
type Execer struct {
	solvers []string
	logger  *log.Logger
}

// New creates a new Execer instance.
func New(solvers []string, logger *log.Logger) *Execer {
	return &Execer{solvers: solvers, logger: logger}
}

// Run runs the registered solvers.
func (e *Execer) Run(args []string) <-chan *Result {
	numSolvers := len(e.solvers)
	eventCh := make(chan *Result, numSolvers)
	ctx := context.Background()

	go func() {
		wg := new(sync.WaitGroup)
		wg.Add(numSolvers)
		for _, solver := range e.solvers {
			go e.runSolver(ctx, wg, solver, args, eventCh)
		}
		wg.Wait()
		close(eventCh)
	}()
	return eventCh
}

func isResult(b []byte) bool {
	m := make(map[string]any)
	if err := json.Unmarshal(b, &m); err != nil {
		return false
	}
	result, ok := m["msg"]
	if !ok {
		return false
	}
	if result != "result" {
		return false
	}
	if _, ok := m["moves"]; !ok {
		return false
	}
	return true
}

var errResultNotFound = errors.New("result not found")

func (e *Execer) runSolver(ctx context.Context, wg *sync.WaitGroup, solver string, args []string, resultCh chan<- *Result) {
	defer wg.Done()

	result := new(Result)
	result.Response, result.Err = func() ([]byte, error) {

		var result []byte
		cmd := exec.CommandContext(ctx, solver, args...)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}
		defer stdout.Close()

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return nil, err
		}
		defer stderr.Close()

		if err := cmd.Start(); err != nil {
			return nil, err
		}

		outScanner := bufio.NewScanner(stdout)
		for outScanner.Scan() {
			b := outScanner.Bytes()
			if isResult(b) {
				result = b
			}
			fmt.Println(string(b))
		}
		if err := outScanner.Err(); err != nil {
			return nil, err
		}

		errScanner := bufio.NewScanner(stderr)
		for errScanner.Scan() {
			fmt.Println(string(outScanner.Bytes()))
		}
		if err := errScanner.Err(); err != nil {
			return nil, err
		}

		if err := cmd.Wait(); err != nil {
			return nil, err
		}

		if result == nil {
			return nil, errResultNotFound
		}
		return result, nil

	}()
	resultCh <- result
}
