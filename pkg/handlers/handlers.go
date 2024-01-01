package handlers

import (
	"github.com/Sunpacker/go-booking-app/pkg/config"
	"github.com/Sunpacker/go-booking-app/pkg/models"
	"github.com/Sunpacker/go-booking-app/pkg/render"
	"net/http"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func CreateNewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func SetNewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	_ = r
	render.RenderTemplate(w, "home", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	_ = r

	stringMap := make(map[string]string)
	stringMap["test"] = "hi again"

	render.RenderTemplate(w, "about", &models.TemplateData{
		StringMap: stringMap,
	})
}
