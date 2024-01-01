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
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.RenderTemplate(w, "home", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")

	stringMap := make(map[string]string)
	stringMap["test"] = remoteIp

	render.RenderTemplate(w, "about", &models.TemplateData{
		StringMap: stringMap,
	})
}
