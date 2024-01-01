package main

import (
	"fmt"
	"github.com/Sunpacker/go-booking-app/pkg/config"
	"github.com/Sunpacker/go-booking-app/pkg/handlers"
	"github.com/Sunpacker/go-booking-app/pkg/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"time"
)

const PORT = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func initSession() {
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Persist = true
	session.Cookie.Secure = app.IsProd

	app.Session = session
}

func initPages() {
	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = templateCache
	render.NewTemplates(&app)
}

func initHandlers() {
	repo := handlers.CreateNewRepo(&app)
	handlers.SetNewHandlers(repo)
}

func main() {
	app.IsProd = false
	app.UseCache = false

	initSession()
	initHandlers()
	initPages()

	serve := &http.Server{
		Addr:    PORT,
		Handler: routes(),
	}

	fmt.Printf("Starting application on port %s\n", PORT)
	err := serve.ListenAndServe()
	log.Fatal(err)
}
