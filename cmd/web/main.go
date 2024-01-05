package main

import (
	"encoding/gob"
	"fmt"
	"github.com/Sunpacker/go-booking-app/internal/config"
	"github.com/Sunpacker/go-booking-app/internal/driver"
	"github.com/Sunpacker/go-booking-app/internal/handlers"
	"github.com/Sunpacker/go-booking-app/internal/helpers"
	"github.com/Sunpacker/go-booking-app/internal/models"
	"github.com/Sunpacker/go-booking-app/internal/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"time"
)

const PORT = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	_, err := run()
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

func run() (*driver.DB, error) {
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Reservation{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})

	app.IsProd = false
	app.UseCache = app.IsProd
	app.InfoLog = log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=booking-app user=postgres password=root")
	if err != nil {
		log.Fatal("cannot connect to database, dying...")
	}

	initSession()
	initHandlers(db)
	err = initPages()
	if err != nil {
		return nil, err
	}
	initHelpers()

	return db, nil
}

func initSession() {
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Persist = true
	session.Cookie.Secure = app.IsProd

	app.Session = session
}

func initHandlers(db *driver.DB) {
	repo := handlers.CreateNewRepo(&app, db)
	handlers.NewHandlers(repo)
}

func initPages() error {
	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}

	app.TemplateCache = templateCache
	render.NewRenderer(&app)

	return nil
}

func initHelpers() {
	helpers.NewHelpers(&app)
}
