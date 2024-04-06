package script

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"github.com/enchik0reo/commandApi/internal/logs"
)

type ScriptExecutor struct {
	log *logs.CustomLog
}

func NewExecutor(log *logs.CustomLog) *ScriptExecutor {
	return &ScriptExecutor{log: log}
}

func (e *ScriptExecutor) RunScript(script string, stop <-chan struct{}) (<-chan string, <-chan error) {
	const op = "script.StartScript"
	out := make(chan string)
	errOut := make(chan error)

	go func() {
		defer close(errOut)
		defer close(out)

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
					e.log.Debug("script stopped manually", e.log.Attr("op", op), e.log.Attr("script", script))

					if err := stdout.Close(); err != nil {
						errOut <- fmt.Errorf("can't close stdout pipe: %s: %v", op, err)
					}
				}
			}
		}()

		for scanner.Scan() {
			out <- scanner.Text()
		}

		if err := cmd.Wait(); err != nil {
			if !strings.Contains(err.Error(), syscall.EPIPE.Error()) {
				errOut <- fmt.Errorf("can't execute script: %v", err)
				return
			} else {
				errOut <- syscall.EPIPE
				return
			}
		}
	}()

	return out, errOut
}
