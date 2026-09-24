package main

import (
	"log"
	"net/http"
	"time"

	"github.com/DesirArman007/social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

// func (app *application) mount() *http.ServeMux {

// 	mux := http.NewServeMux()

// 	mux.HandleFunc("GET /v1/health", app.healthCheckHandler)
// 	return mux

// }

func (app *application) mount() http.Handler {
	route := chi.NewRouter()

	route.Use(middleware.RequestID)
	route.Use(middleware.RealIP)
	route.Use(middleware.Logger)
	route.Use(middleware.Recoverer)

	// set a timeout on the request context (Ctx) , that will singal thorugh
	// ctx.Done() that the req has timedout and further processing should be stopped
	route.Use(middleware.Timeout(60 * time.Second))

	route.Route("/v1", func(route chi.Router) {
		route.Get("/health", app.healthCheckHandler)
	})

	return route
}

func (app *application) run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at %s", app.config.addr)

	return srv.ListenAndServe()

}
