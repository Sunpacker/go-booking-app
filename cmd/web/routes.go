package main

import (
	"github.com/Sunpacker/go-booking-app/pkg/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func initMiddlewares(mux *chi.Mux) {
	mux.Use(middleware.AllowContentType())
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
}

func initStaticFilesDir(mux *chi.Mux) {
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
}

func initPageRoutes(mux *chi.Mux) {
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
}

func routes() http.Handler {
	mux := chi.NewRouter()

	initMiddlewares(mux)
	initStaticFilesDir(mux)
	initPageRoutes(mux)

	return mux
}
