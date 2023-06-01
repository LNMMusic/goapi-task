package application

import (
	"api/cmd/rest/handlers"
	"api/cmd/rest/middlewares/logger"
	"api/internal/task"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// -----------------------------------------------------------------------------
func NewConfigDefault() *Config {
	return &Config{}
}

// Config is an struct that contains all the configuration of the application.
type Config struct {
	// ...
}


// -----------------------------------------------------------------------------
func NewApp(config *Config, router chi.Router) *App {
	return &App{config: config, router: router}
}
// App is an struct that contains all the dependencies of the application.
type App struct {
	// config: represents the configuration of the application.
	config *Config
	// router: represents the| router of the application.
	router chi.Router
}

func (a *App) Dependencies() (err error) {
	// initialize dependencies (based on config)
	db := []*task.Task{}
	vl := task.NewValidatorLocal()
	st := task.NewStorageLocal(db, vl)

	ct := handlers.NewTaskController(st)

	// register routes
	// -> middlewares: handler#1 -> (http.HandlerFunc) middleware #1 -> (http.Handler) middleware #2 -> ... -> serveHTTP()
	a.router.Use(middleware.Recoverer)
	a.router.Use(logger.LoggerDefault)

	// -> handlers
	a.router.Get("/ping", handlers.Health())

	a.router.Route("/tasks", func(r chi.Router) {
		// Get a task
		r.Get("/{id}", ct.Get())
		// Create a task
		r.Post("/", ct.Create())
	})

	return
}

// Run starts the application.
func (a *App) Run() (err error) {
	// start the application
	err = http.ListenAndServe(":8080", a.router)
	return
}