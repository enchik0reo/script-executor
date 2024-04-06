package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/enchik0reo/commandApi/internal/config"
	"github.com/enchik0reo/commandApi/internal/logs"
	"github.com/enchik0reo/commandApi/internal/server/handler"
	"github.com/enchik0reo/commandApi/internal/server/server"
	"github.com/enchik0reo/commandApi/internal/services/commander"
	"github.com/enchik0reo/commandApi/internal/services/script"
	"github.com/enchik0reo/commandApi/internal/storage"
)

type App struct {
	cfg *config.Config
	log *logs.CustomLog
	db  *sql.DB
	cmd *commander.Commander
	srv *server.Server
}

func New() *App {
	a := &App{}
	var err error

	a.cfg = config.MustLoad()

	a.log = logs.NewLogger(a.cfg.Env)

	a.db, err = connectionAttemptToDB(a.cfg.Storage)
	if err != nil {
		a.log.Error("Failed to connect to db", a.log.Attr("error", err))
		os.Exit(1)
	}

	cS := storage.NewCommandStorage(a.db)

	e := script.NewExecutor(a.log)

	a.cmd = commander.NewCommander(a.log, cS, e)

	handler, err := handler.New(a.cmd, a.cfg.Frontend.Domains, a.cfg.Server.Timeout, a.log)
	if err != nil {
		a.log.Error("Failed to create handler", a.log.Attr("error", err))
		os.Exit(1)
	}

	a.srv = server.New(handler, &a.cfg.Server, a.log)

	return a
}

func (a *App) MustRun() {
	a.log.Info("Starting command executor service", "env", a.cfg.Env)

	go func() {
		if err := a.srv.Start(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				a.log.Error("Failed ower working api service", a.log.Attr("error", err))
				os.Exit(1)
			}
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	a.mustStop()
}

func (a *App) mustStop() {
	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.CtxTimeout)
	defer cancel()

	if err := a.cmd.StopAllRunningScripts(ctx); err != nil {
		a.log.Error("Stoping running commands", a.log.Attr("error", err))
	}

	if err := a.srv.Stop(ctx); err != nil {
		a.log.Error("Closing connection to api server", a.log.Attr("error", err))
	}

	if err := a.db.Close(); err != nil {
		a.log.Error("Closing connection to command storage", a.log.Attr("error", err))
	}

	a.log.Info("Command executor  stoped gracefully")
}

func connectionAttemptToDB(psql config.Postgres) (*sql.DB, error) {
	var err error
	var db *sql.DB

	for i := 1; i <= 5; i++ {
		db, err = storage.Connect(psql)
		if err != nil {
			time.Sleep(time.Duration(i) * time.Second)
		} else {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("after retries: %w", err)
	}

	return db, nil
}
