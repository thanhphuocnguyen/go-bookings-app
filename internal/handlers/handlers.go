package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
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

const layout = "2006-01-02"

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

var Repo *Repository

func SetRepository(r *Repository) {
	Repo = r
}

func InitializeRepository(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbRepo.InitPGRepository(a, db.SQL),
	}
}

func InitializeTestingRepository(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbRepo.InitTestingRepository(a, nil),
	}
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	m.DB.AllUsers()
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "searchAvailability.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	startDate, err := time.Parse(layout, r.Form.Get("start_date"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, r.Form.Get("end_date"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityInRange(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No available rooms")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
	}

	data := make(map[string]interface{})

	data["rooms"] = rooms
	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "chooseRoom.page.tmpl", &models.TemplateData{
		Data: data})
}

func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{}
	statusCode := http.StatusOK

	err := r.ParseForm()
	if err != nil {
		resp.OK = false
		statusCode = http.StatusBadRequest
		resp.Message = "Cannot parse form"
	} else {
		startDate, errParseSD := time.Parse(layout, r.Form.Get("start"))
		endDate, errParsedED := time.Parse(layout, r.Form.Get("end"))
		roomId, errParsedRoomId := strconv.Atoi(r.Form.Get("room_id"))
		if errParsedED != nil || errParseSD != nil {
			resp.OK = false
			statusCode = http.StatusBadRequest
			resp.Message = "Cannot parse dates"
		} else if errParsedRoomId != nil {
			resp.OK = false
			statusCode = http.StatusBadRequest
			resp.Message = "Cannot parse room id"
		} else {
			available, err := m.DB.CheckIfRoomAvailableByDate(roomId, startDate, endDate)
			if err != nil {
				resp.OK = false
				statusCode = http.StatusInternalServerError
				resp.Message = "Error checking room availability"
			} else {
				resp.OK = available
				resp.StartDate = startDate.Format(layout)
				resp.EndDate = endDate.Format(layout)
			}
		}
	}
	out, _ := json.MarshalIndent(resp, "", "  ")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(out)
}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot parse form")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	f := forms.New(r.PostForm)
	f.Required("email", "password")
	f.IsEmail("email")

	if !f.Valid() {
		strMap := make(map[string]string)
		strMap["email"] = email
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form:      f,
			StringMap: strMap,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Invalid login")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) ShowRegistration(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "registration.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostRegistration(w http.ResponseWriter, r *http.Request) {

}
