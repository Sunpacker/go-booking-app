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

	serve := &http.Server{
		Addr:    PORT,
		Handler: routes(),
	}

	fmt.Printf("Starting application on port %s\n", PORT)
	err = serve.ListenAndServe()
	log.Fatal(err)
}
