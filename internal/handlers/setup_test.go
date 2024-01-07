package handlers

import (
	"encoding/gob"
	"fmt"
	"github.com/Sunpacker/go-booking-app/internal/config"
	"github.com/Sunpacker/go-booking-app/internal/driver"
	"github.com/Sunpacker/go-booking-app/internal/models"
	"github.com/Sunpacker/go-booking-app/internal/render"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var app config.AppConfig
var session *scs.SessionManager
var functions = template.FuncMap{
	"formatDate": render.FormatDate,
}

func TestMain(m *testing.M) {
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Reservation{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})

	app.IsProd = false
	app.UseCache = true
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
		log.Fatal("cannot create template cache")
	}

	os.Exit(m.Run())
}

func getRoutes() http.Handler {
	mux := chi.NewRouter()

	initMiddlewares(mux)
	initStaticFilesDir(mux)
	initPageRoutes(mux)

	return mux
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
	_ = db
	repo := NewTestRepo(&app)
	NewHandlers(repo)
}

func initPages() error {
	templateCache, err := CreateTestTemplateCache()
	if err != nil {
		return err
	}

	app.TemplateCache = templateCache
	render.NewRenderer(&app)

	return nil
}
func getTemplateFilepathPattern(pattern string) string {
	return fmt.Sprintf("./../../templates/%s.tmpl", pattern)
}
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	templateCache := map[string]*template.Template{}

	pagePattern := getTemplateFilepathPattern("*.page")
	pages, err := filepath.Glob(pagePattern)
	if err != nil {
		return templateCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		templateSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		layoutPattern := getTemplateFilepathPattern("*.layout")
		matches, err := filepath.Glob(layoutPattern)
		if err != nil {
			return templateCache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob(layoutPattern)
			if err != nil {
				return templateCache, err
			}
		}

		pageName := strings.Split(name, ".")[0]
		templateCache[pageName] = templateSet
	}

	return templateCache, nil
}

func initMiddlewares(mux *chi.Mux) {
	mux.Use(middleware.Recoverer)
	mux.Use(SessionLoad)
}
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func initStaticFilesDir(mux *chi.Mux) {
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
}

func initPageRoutes(mux *chi.Mux) {
	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Get("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)

	mux.Get("/contact", Repo.Contact)
}
