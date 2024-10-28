package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/thanhphuocnguyen/go-bookings-app/internal/config"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/driver"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/forms"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/helpers"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/models"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/render"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/repository"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/repository/dbRepo"
)

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

var Repo *Repository

func InitRepo(r *Repository) {
	Repo = r
}

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbRepo.InitPGRepository(a, db.SQL),
	}
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	m.DB.AllUsers()
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.RenderTemplate(w, r, "makeReservation.page.tmpl", &models.TemplateData{
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

	startDate, err := time.Parse("2006-01-02", "2021-01-01")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse("2006-01-02", "2021-01-02")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		UserId:    1,
		RoomId:    1,
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
	}

	f := forms.New(r.PostForm)

	f.Required("first_name", "last_name", "email")
	f.MinLength("first_name", 3, r)
	f.MinLength("last_name", 3, r)
	f.IsEmail("email")

	if !f.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.RenderTemplate(w, r, "makeReservation.page.tmpl", &models.TemplateData{
			Form: f,
			Data: data,
		})
		return
	}

	newResId, err := m.DB.InsertReservation(&reservation)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	roomRestriction := models.RoomRestriction{
		RoomId:        1,
		StartDate:     startDate,
		EndDate:       endDate,
		RestrictionId: 1,
		ReservationId: newResId,
	}

	m.DB.InsertRoomRestriction(&roomRestriction)
	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Cannot get reservation from session")
		m.App.Session.Put(r.Context(), "error", "Cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.RenderTemplate(w, r, "reservationSummary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "searchAvailability.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {

}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		true,
		"Available!",
	}

	out, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		helpers.ServerError(w, err)
	}
	w.Header().Set("Content-Type", "application/json")
	// set status code to 200
	w.WriteHeader(http.StatusOK)
	log.Println(string(out))
	w.Write(out)
}
