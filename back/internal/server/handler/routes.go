package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/enchik0reo/commandApi/internal/logs"
	"github.com/enchik0reo/commandApi/internal/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Commander interface {
	CreateNewCommand(context.Context, string) (int64, error)
	GetCommandList(context.Context, int64) ([]models.Command, error)
	GetOneCommandDescription(context.Context, int64) (*models.Command, error)
	StopCommand(context.Context, int64) (int64, error)
}

type CustomRouter struct {
	*chi.Mux
	cmdr    Commander
	timeout time.Duration
	log     *logs.CustomLog
}

// New returns new handler ...
func New(cmdr Commander, domains []string, timeout time.Duration, log *logs.CustomLog) (http.Handler, error) {
	r := CustomRouter{chi.NewRouter(), cmdr, timeout, log}

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(loggerMw(log))
	r.Use(corsSettings(domains))

	r.Post("/create", r.create())
	r.Post("/create/upload", r.createUpload())
	r.Get("/list", r.commands())
	r.Get("/cmd", r.command())
	r.Put("/stop", r.stopCommand())

	return r, nil
}
