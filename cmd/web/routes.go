package main

import (
	"github.com/Sunpacker/go-booking-app/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	initMiddlewares(mux)
	initStaticFilesDir(mux)
	initPageRoutes(mux)

	return mux
}

func initMiddlewares(mux *chi.Mux) {
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
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)

	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Get("/search-availability-json", handlers.Repo.AvailabilityJSON)

	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)

	mux.Get("/contact", handlers.Repo.Contact)
}
