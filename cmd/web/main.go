package main

import (
	"fmt"
	"github.com/Sunpacker/go-booking-app/pkg/config"
	"github.com/Sunpacker/go-booking-app/pkg/handlers"
	"github.com/Sunpacker/go-booking-app/pkg/render"
	"log"
	"net/http"
)

const PORT = ":8080"

func main() {
	var app config.AppConfig

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = templateCache
	app.UseCache = false

	repo := handlers.CreateNewRepo(&app)
	handlers.SetNewHandlers(repo)

	render.NewTemplates(&app)

	http.HandleFunc("/", handlers.Repo.Home)
	http.HandleFunc("/about", handlers.Repo.About)

	fmt.Printf("Starting application on port %s\n", PORT)
	_ = http.ListenAndServe(PORT, nil)
}
