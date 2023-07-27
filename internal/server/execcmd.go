package server

import (
	"context"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/go-ricrob/exec/task"
)

const maxCh = 100

type result struct {
	Solver    string     `json:"solver"`
	StartTime time.Time  `json:"startTime"`
	EndTime   time.Time  `json:"endTime"`
	Task      *task.Task `json:"task"`
	Moves     task.Moves `json:"moves"`
	Err       error      `json:"err"`
}

type execCmd struct {
	solvers []string
	l       *log.Logger
}

func newExecCmd(solvers []string, l *log.Logger) *execCmd {
	return &execCmd{solvers: solvers, l: l}
}

func (e *execCmd) execute(task *task.Task) <-chan any {
	e.l.Printf("add task %v", task)

	eventCh := make(chan any, maxCh)
	ctx := context.Background()

	go func() {
		wg := new(sync.WaitGroup)
		wg.Add(len(e.solvers))
		for _, solver := range e.solvers {
			go e.executeSolver(ctx, wg, solver, task, eventCh)
		}
		wg.Wait()
		close(eventCh)
	}()
	return eventCh
}

func (e *execCmd) executeSolver(ctx context.Context, wg *sync.WaitGroup, solver string, task *task.Task, eventCh chan<- any) {
	defer wg.Done()

	result := &result{
		Solver:    solver,
		StartTime: time.Now(),
		Task:      task,
	}

	result.Err = func() error {

		cmd := exec.CommandContext(ctx, solver, task.Args.CmdArgs()...)

		if err := cmd.Start(); err != nil {
			return err
		}

		if err := cmd.Wait(); err != nil {
			return err
		}

		return nil

	}()

	//...

	eventCh <- result
}
