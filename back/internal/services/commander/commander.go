package commander

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/enchik0reo/commandApi/internal/logs"
	"github.com/enchik0reo/commandApi/internal/models"
	"github.com/enchik0reo/commandApi/internal/services"
)

type Storager interface {
	CreateNew(context.Context, string) (int64, error)
	GetList(context.Context, int64) ([]models.Command, error)
	GetOne(context.Context, int64) (*models.Command, error)
	StopOne(context.Context, int64) (int64, error)
	SaveOutput(context.Context, int64, string) (int64, error)
}

type Executor interface {
	RunScript(string, string, <-chan struct{}) (<-chan string, <-chan error)
}

const (
	contextDuration = 3 * time.Second
	maxScriptLenght = 27
)

type Commander struct {
	cmdStorage Storager
	exec       Executor

	log       *logs.CustomLog
	stopChans map[int64]chan struct{}
	mu        *sync.RWMutex
}

// NewCommander creates a new instance of Commander ...
func NewCommander(l *logs.CustomLog, s Storager, e Executor) *Commander {
	c := &Commander{
		log:        l,
		cmdStorage: s,
		exec:       e,
		stopChans:  make(map[int64]chan struct{}),
		mu:         &sync.RWMutex{},
	}

	return c
}

// CreateNewCommand starts new script.
// It creates new record in storage and runs the script in new gorutine ...
func (c *Commander) CreateNewCommand(ctx context.Context, script string) (int64, error) {
	const op = "commander.CreateNewCommand"

	sName, err := scriptName(script)
	if err != nil {
		return -1, fmt.Errorf("can't make script name: %s: %v", op, err)
	}

	id, err := c.cmdStorage.CreateNew(ctx, sName)
	if err != nil {
		return -1, fmt.Errorf("can't create new command in storage: %s: %v", op, err)
	}

	stopCh := make(chan struct{})

	c.mu.Lock()
	c.stopChans[id] = stopCh
	c.mu.Unlock()

	resCh, errCh := c.exec.RunScript(script, sName, stopCh)

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
					if errors.Is(err, services.ErrStoppedManually) {
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

// GetCommandList returns the list of command with limit from storage ...
func (c *Commander) GetCommandList(ctx context.Context, limit int64) ([]models.Command, error) {
	const op = "commander.GetCommandList"
	cmds, err := c.cmdStorage.GetList(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("can't get list if command: %s: %v", op, err)
	}

	return cmds, nil
}

// GetOneCommandDescription returns the command's info from storage ...
func (c *Commander) GetOneCommandDescription(ctx context.Context, id int64) (*models.Command, error) {
	const op = "commander.GetCommandDescription"
	cmd, err := c.cmdStorage.GetOne(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("can't get command description on id: %d: %s: %v", id, op, err)
	}

	return cmd, nil
}

// StopCommand stops the command by id.
// It updates record in storage ...
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

// StopAllRunningScripts stops all running commands ...
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

// scriptName returns correct name of script ...
func scriptName(script string) (string, error) {
	sName := strings.ReplaceAll(script, "\n", " ")
	res := []rune(sName)

	if len(res) > maxScriptLenght {
		res = res[:maxScriptLenght]
		res = append(res, []rune("...")...)
	}

	return string(res), nil
}
