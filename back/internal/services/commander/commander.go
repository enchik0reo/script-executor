package commander

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/enchik0reo/commandApi/internal/logs"
	"github.com/enchik0reo/commandApi/internal/models"
	"github.com/enchik0reo/commandApi/internal/services"
)

type Sotrager interface {
	CreateNew(context.Context, string) (int64, error)
	GetList(context.Context, int64) ([]models.Command, error)
	GetOne(context.Context, int64) (*models.Command, error)
	StopOne(context.Context, int64) (int64, error)
	SaveOutput(context.Context, int64, string) (int64, error)
}

type Executor interface {
	RunScript(string, <-chan struct{}) (<-chan string, <-chan error)
}

const (
	contextDuration = 3 * time.Second
	minScriptLenght = 11
	maxScriptLenght = 38
)

type Commander struct {
	cmdStorage Sotrager
	exec       Executor

	log       *logs.CustomLog
	stopChans map[int64]chan struct{}
	mu        *sync.RWMutex
}

func NewCommander(l *logs.CustomLog, s Sotrager, e Executor) *Commander {
	c := &Commander{
		log:        l,
		cmdStorage: s,
		exec:       e,
		stopChans:  make(map[int64]chan struct{}),
		mu:         &sync.RWMutex{},
	}

	return c
}

func (c *Commander) CreateNewCommmand(ctx context.Context, script string) (int64, error) {
	const op = "commander.CreateNewCommand"

	scriptName := validateScript(script)

	id, err := c.cmdStorage.CreateNew(ctx, scriptName)
	if err != nil {
		return -1, fmt.Errorf("can't create new command in storage: %s: %v", op, err)
	}

	stopCh := make(chan struct{})

	c.mu.Lock()
	c.stopChans[id] = stopCh
	c.mu.Unlock()

	resCh, errCh := c.exec.RunScript(script, stopCh)

	go func() {
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), contextDuration)

			if _, err := c.cmdStorage.StopOne(ctx, id); err != nil {
				c.log.Error("can't save output in storage", c.log.Attr("op", op), c.log.Attr("error", err))
			}

			cancel()

			c.mu.Lock()
			delete(c.stopChans, id)
			c.mu.Unlock()

			close(stopCh)
		}()

		for {
			select {
			case res, open := <-resCh:
				if open {
					ctx, cancel := context.WithTimeout(context.Background(), contextDuration)

					if _, err := c.cmdStorage.SaveOutput(ctx, id, res); err != nil {
						c.log.Error("can't save output in storage", c.log.Attr("op", op), c.log.Attr("error", err))
					}

					cancel()
				} else {
					return
				}
			case err, open := <-errCh:
				if open {
					if errors.Is(err, syscall.EPIPE) {
						ctx, cancel := context.WithTimeout(context.Background(), contextDuration)

						if _, errOut := c.cmdStorage.SaveOutput(ctx, id, "Execution was interrupted"); errOut != nil {
							c.log.Error("can't save info in storage", c.log.Attr("op", op), c.log.Attr("error", errOut))
						}

						cancel()
						return
					} else {
						ctx, cancel := context.WithTimeout(context.Background(), contextDuration)

						if _, errOut := c.cmdStorage.SaveOutput(ctx,
							id, fmt.Sprintf("Stopped with error: %s", err.Error())); errOut != nil {
							c.log.Error("can't save error in storage", c.log.Attr("op", op), c.log.Attr("error", errOut))
						}

						cancel()

						c.log.Info("can't execute sctipt", c.log.Attr("op", op), c.log.Attr("error", err))
						return
					}
				} else {
					return
				}
			}
		}
	}()

	return id, nil
}

func (c *Commander) GetCommandList(ctx context.Context, limit int64) ([]models.Command, error) {
	const op = "commander.GetCommandList"
	cmds, err := c.cmdStorage.GetList(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("can't get list if command: %s: %v", op, err)
	}

	return cmds, nil
}

func (c *Commander) GetOneCommandDescription(ctx context.Context, id int64) (*models.Command, error) {
	const op = "commander.GetCommandDescription"
	cmd, err := c.cmdStorage.GetOne(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("can't get command description on id: %d: %s: %v", id, op, err)
	}

	return cmd, nil
}

func (c *Commander) StopCommand(ctx context.Context, id int64) (int64, error) {
	const op = "commander.StopCommand"
	var res int64
	var err error

	c.mu.RLock()
	ch, ok := c.stopChans[id]
	c.mu.RUnlock()
	if ok {
		ch <- struct{}{}

		c.mu.Lock()
		delete(c.stopChans, id)
		c.mu.Unlock()

		if res, err = c.cmdStorage.StopOne(ctx, id); err != nil {
			c.log.Error("can't save output in storage", c.log.Attr("op", op), c.log.Attr("error", err))
		}
	} else {
		return 0, services.ErrNoExecutingCommand
	}

	return res, nil
}

func (c *Commander) StopAllRunningScripts(ctx context.Context) error {
	const op = "commander.StopAllRunningScripts"
	var resErr error

	for id, ch := range c.stopChans {
		ch <- struct{}{}

		c.mu.Lock()
		delete(c.stopChans, id)
		c.mu.Unlock()

		if _, err := c.cmdStorage.StopOne(ctx, id); err != nil {
			c.log.Error("can't save output in storage", c.log.Attr("op", op), c.log.Attr("error", err))
			resErr = errors.Join(err)
		}
	}

	return resErr
}

func validateScript(script string) string {
	scriptName := strings.ReplaceAll(script, "\n", " ")

	if len(scriptName) > maxScriptLenght {
		scriptName = scriptName[minScriptLenght:maxScriptLenght] + "..."
	} else if len(scriptName) > minScriptLenght {
		scriptName = scriptName[minScriptLenght:]
	}

	return scriptName
}
