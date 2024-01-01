package main

import (
	"github.com/Sunpacker/go-booking-app/pkg/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.AllowContentType())
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	return mux
}