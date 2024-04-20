package handler

import (
	"context"
	"net/http"
	"time"

	_ "github.com/enchik0reo/commandApi/docs"
	"github.com/enchik0reo/commandApi/internal/logs"
	"github.com/enchik0reo/commandApi/internal/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger" // swagger embed files
)

//go:generate mockgen -destination=mocks/handler.go -package=mocks -source=routes.go

//go:generate go run github.com/vektra/mockery/v2@v2.42.2 --name=Commander
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
func New(cmdr Commander, domains []string, timeout time.Duration, log *logs.CustomLog) http.Handler {
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

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8008/swagger/doc.json"),
	))

	return r
}
