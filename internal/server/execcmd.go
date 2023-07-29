package server

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

const maxCh = 100

type event struct {
	result []byte
	err    error
}

type execCmd struct {
	solvers []string
	logger  *log.Logger
}

func newExecCmd(solvers []string, logger *log.Logger) *execCmd {
	return &execCmd{solvers: solvers, logger: logger}
}

func (e *execCmd) execute(args []string) <-chan *event {
	eventCh := make(chan *event, maxCh)
	ctx := context.Background()

	go func() {
		wg := new(sync.WaitGroup)
		wg.Add(len(e.solvers))
		for _, solver := range e.solvers {
			go e.executeSolver(ctx, wg, solver, args, eventCh)
		}
		wg.Wait()
		close(eventCh)
	}()
	return eventCh
}

func (e *execCmd) isResult(b []byte) bool {
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

func (e *execCmd) executeSolver(ctx context.Context, wg *sync.WaitGroup, solver string, args []string, eventCh chan<- *event) {
	defer wg.Done()

	event := new(event)
	event.result, event.err = func() ([]byte, error) {

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
			if e.isResult(b) {
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
	eventCh <- event
}
