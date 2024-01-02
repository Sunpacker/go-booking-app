package main

import (
	"encoding/gob"
	"fmt"
	"github.com/Sunpacker/go-booking-app/internal/config"
	"github.com/Sunpacker/go-booking-app/internal/handlers"
	"github.com/Sunpacker/go-booking-app/internal/models"
	"github.com/Sunpacker/go-booking-app/internal/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"time"
)

const PORT = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	serve := &http.Server{
		Addr:    PORT,
		Handler: routes(),
	}

	fmt.Printf("Starting application on port %s\n", PORT)
	err = serve.ListenAndServe()
	log.Fatal(err)
}

func run() error {
	gob.Register(models.Reservation{})
	app.IsProd = false
	app.UseCache = app.IsProd

	initSession()

	initHandlers()

	err := initPages()
	if err != nil {
		return err
	}

	return nil
}

func initSession() {
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Persist = true
	session.Cookie.Secure = app.IsProd

	app.Session = session
}

func initHandlers() {
	repo := handlers.CreateNewRepo(&app)
	handlers.SetNewHandlers(repo)
}

func initPages() error {
	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}

	app.TemplateCache = templateCache
	render.NewTemplates(&app)

	return nil
}
