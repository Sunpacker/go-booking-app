package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Sunpacker/go-booking-app/internal/config"
	"github.com/Sunpacker/go-booking-app/internal/forms"
	"github.com/Sunpacker/go-booking-app/internal/helpers"
	"github.com/Sunpacker/go-booking-app/internal/models"
	"github.com/Sunpacker/go-booking-app/internal/render"
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

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	_ = render.RenderTemplate(w, r, "home", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	_ = render.RenderTemplate(w, r, "about", &models.TemplateData{})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation

	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	_ = render.RenderTemplate(w, r, "make-reservation", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		_ = render.RenderTemplate(w, r, "make-reservation", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	_ = render.RenderTemplate(w, r, "generals", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	_ = render.RenderTemplate(w, r, "majors", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	_ = render.RenderTemplate(w, r, "search-availability", &models.TemplateData{})
}
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	_, _ = w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	_ = r

	resp := jsonResponse{
		Ok:      true,
		Message: "Available!",
	}

	out, err := json.Marshal(&resp)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(out)
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	_ = render.RenderTemplate(w, r, "contact", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	const reservationSession = "reservation"

	reservation, ok := m.App.Session.Get(r.Context(), reservationSession).(models.Reservation)
	if !ok {
		msg := "cannot get error from session"
		m.App.ErrorLog.Println(msg)
		m.App.Session.Put(r.Context(), "error", msg)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), reservationSession)

	data := make(map[string]interface{})
	data[reservationSession] = reservation

	_ = render.RenderTemplate(w, r, "reservation-summary", &models.TemplateData{
		Data: data,
	})
}
