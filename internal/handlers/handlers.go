package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Sunpacker/go-booking-app/internal/config"
	"github.com/Sunpacker/go-booking-app/internal/driver"
	"github.com/Sunpacker/go-booking-app/internal/forms"
	"github.com/Sunpacker/go-booking-app/internal/helpers"
	"github.com/Sunpacker/go-booking-app/internal/models"
	"github.com/Sunpacker/go-booking-app/internal/render"
	"github.com/Sunpacker/go-booking-app/internal/repository"
	"github.com/Sunpacker/go-booking-app/internal/repository/dbrepo"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func CreateNewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestRepo(a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	_ = render.Template(w, r, "home", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	_ = render.Template(w, r, "about", &models.TemplateData{})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomByID(reservation.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.Room.RoomName = room.RoomName

	startDate := reservation.StartDate.Format("2006-01-02")
	endDate := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = startDate
	stringMap["end_date"] = endDate

	data := make(map[string]interface{})
	data["reservation"] = reservation

	_ = render.Template(w, r, "make-reservation", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	// 2020-01-01 -- 01/02 03:04:05PM '06 -0700

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		http.Error(w, "my own error message", http.StatusSeeOther)
		_ = render.Template(w, r, "make-reservation", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation into database!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert room restriction!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	_ = render.Template(w, r, "generals", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	_ = render.Template(w, r, "majors", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	_ = render.Template(w, r, "search-availability", &models.TemplateData{})
}
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	// 01/02 03:04:05PM '06 -0700
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get availability for rooms")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No available room")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	reservation := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	_ = render.Template(w, r, "choose-room", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	Ok        bool   `json:"ok"`
	Message   string `json:"message"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	RoomID    string `json:"room_id"`
}

func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		resp := jsonResponse{
			Ok:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		// got a database error, so return appropriate json
		resp := jsonResponse{
			Ok:      false,
			Message: "Error querying database",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(out)
		return
	}
	resp := jsonResponse{
		Ok:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	// I removed the error check, since we handle all aspects of
	// the json right here
	out, _ := json.MarshalIndent(resp, "", "     ")

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(out)
}
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(exploded[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	var reservation models.Reservation
	queryParams := r.URL.Query()

	roomID, _ := strconv.Atoi(queryParams.Get("id"))
	reservation.RoomID = roomID

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, queryParams.Get("s"))
	endDate, _ := time.Parse(layout, queryParams.Get("e"))
	reservation.StartDate = startDate
	reservation.EndDate = endDate

	room, err := m.DB.GetRoomByID(reservation.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get room from db!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	reservation.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	_ = render.Template(w, r, "contact", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	_ = render.Template(w, r, "reservation-summary", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = render.Template(w, r, "login", &models.TemplateData{
		Form: forms.New(nil),
	})
}
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		_ = render.Template(w, r, "login", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// ADMIN ROUTES //

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	_ = render.Template(w, r, "admin-dashboard", &models.TemplateData{})
}

func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	_ = render.Template(w, r, "admin-new-reservations", &models.TemplateData{
		Data: data,
	})
}
func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	_ = render.Template(w, r, "admin-all-reservations", &models.TemplateData{
		Data: data,
	})
}
func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	_ = render.Template(w, r, "admin-reservations-calendar", &models.TemplateData{})
}

func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	splitted := strings.Split(r.RequestURI, "/")
	src := splitted[3]
	id, err := strconv.Atoi(splitted[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	stringMap := make(map[string]string)
	stringMap["src"] = src

	reservation, err := m.DB.GetReservationByID(id)

	data := make(map[string]interface{})
	data["reservation"] = reservation

	_ = render.Template(w, r, "admin-reservations-show", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil),
	})
}
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_ = m.DB.UpdateProcessedForReservation(id, 1)
	m.App.Session.Put(r.Context(), "flash", "Reservation marked as processed")

	src := chi.URLParam(r, "src")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}
func (m *Repository) AdminPostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	splitted := strings.Split(r.RequestURI, "/")
	src := splitted[3]
	id, err := strconv.Atoi(splitted[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	stringMap := make(map[string]string)
	stringMap["src"] = src

	reservation, err := m.DB.GetReservationByID(id)
	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	err = m.DB.UpdateReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_ = m.DB.DeleteReservation(id)
	m.App.Session.Put(r.Context(), "flash", "Reservation deleted")

	src := chi.URLParam(r, "src")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}
