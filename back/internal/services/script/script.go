package script

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync/atomic"

	"github.com/enchik0reo/commandApi/internal/logs"
	"github.com/enchik0reo/commandApi/internal/services"
)

type Executor struct {
	log *logs.CustomLog
}

// NewExecutor creates a new instance of Executor ...
func NewExecutor(log *logs.CustomLog) *Executor {
	return &Executor{log: log}
}

// RunScript executing script.
// It returns channels for use in new gorutine ...
func (e *Executor) RunScript(script, scriptName string, stop <-chan struct{}) (<-chan string, <-chan error) {
	const op = "script.StartScript"
	var manualStopFlag int32 = 0

	out := make(chan string)
	errOut := make(chan error)

	go func() {
		defer func() {
			close(errOut)
			close(out)
		}()

		cmd := exec.Command("/bin/bash", "-c", script)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			errOut <- fmt.Errorf("can't do stdout pipe: %s: %v", op, err)
			return
		}

		if err := cmd.Start(); err != nil {
			errOut <- fmt.Errorf("can't start script: %s: %v", op, err)
			return
		}

		scanner := bufio.NewScanner(stdout)

		go func() {
			for v := range stop {
				if v == struct{}{} {
					e.log.Debug("script stopped manually", e.log.Attr("op", op), e.log.Attr("script", scriptName))

					if err := stdout.Close(); err != nil {
						errOut <- fmt.Errorf("can't close stdout pipe: %s: %v", op, err)
					}

					atomic.AddInt32(&manualStopFlag, 1)

					if err := cmd.Process.Kill(); err != nil {
						errOut <- fmt.Errorf("can't kill process: %s: %v", op, err)
					}
				}
			}
		}()

		for scanner.Scan() {
			out <- scanner.Text()
		}

		process, err := cmd.Process.Wait()
		if err != nil {
			errOut <- fmt.Errorf("can't wait process: %v", err)
		}

		success := process.Success()
		if !success {
			if manualStopFlag == 1 {
				errOut <- services.ErrStoppedManually
			} else {
				errOut <- fmt.Errorf("can't execute script")
			}
		}
	}()

	return out, errOut
}
