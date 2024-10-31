package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
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
	render.RenderTmpl(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTmpl(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTmpl(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomById(reservation.RoomId)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot get room from database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.Room = room
	emptyReservation.RoomId = reservation.RoomId
	startDate := reservation.StartDate.Format(layout)
	endDate := reservation.EndDate.Format(layout)

	strMap := make(map[string]string)
	strMap["start_date"] = startDate
	strMap["end_date"] = endDate
	strMap["room_name"] = reservation.Room.Name
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTmpl(w, r, "makeReservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: strMap,
	})
}
func (m *Repository) CreateReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.UserId = 1
	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	f := forms.New(r.PostForm)

	f.Required("first_name", "last_name", "email", "phone")
	f.MinLength("first_name", 3, r)
	f.MinLength("last_name", 3, r)
	f.MinLength("phone", 10, r)
	f.IsEmail("email")

	if !f.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		http.Error(w, "Form is not valid", http.StatusSeeOther)
		render.RenderTmpl(w, r, "makeReservation.page.tmpl", &models.TemplateData{
			Form: f,
			Data: data,
		})
		return
	}

	newResId, err := m.DB.InsertReservation(&reservation)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot insert reservation into database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomRestriction := models.RoomRestriction{
		RoomId:        reservation.RoomId,
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RestrictionId: 1,
		ReservationId: newResId,
	}

	htmlMEssage := fmt.Sprintf(`
		<strong>Reservation Confirmation</strong><br>
		Dear %s,<br>
		This is to confirm your reservation from %s to %s.
	`, reservation.FirstName, reservation.StartDate.Format(layout), reservation.EndDate.Format(layout))

	msg := models.MailData{
		To:       "john@doe.ca",
		From:     "universal@booking.com",
		Subject:  "Reservation Confirmation",
		Content:  htmlMEssage,
		Template: "basic.html",
	}

	m.App.MailChan <- msg

	err = m.DB.InsertRoomRestriction(&roomRestriction)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot insert room restriction into database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

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

	strMap := make(map[string]string)
	strMap["start_date"] = reservation.StartDate.Format(layout)
	strMap["end_date"] = reservation.EndDate.Format(layout)
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.RenderTmpl(w, r, "reservationSummary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: strMap,
	})
}

func (m *Repository) Room(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	m.App.InfoLog.Println(id)
	room, err := m.DB.GetRoomBySlug(id)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot get room from database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if room.ID == 0 {
		helpers.ClientError(w, http.StatusNotFound)
		m.App.Session.Put(r.Context(), "error", "Cannot get room from database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data := make(map[string]interface{})
	data["room"] = room

	render.RenderTmpl(w, r, "roomDetails.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTmpl(w, r, "searchAvailability.page.tmpl", &models.TemplateData{
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

	render.RenderTmpl(w, r, "chooseRoom.page.tmpl", &models.TemplateData{
		Data: data})
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot parse room id")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	reservation.UserId = 1
	reservation.RoomId = roomId

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot parse room id")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	startDate, err := time.Parse(layout, r.URL.Query().Get("start_date"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, r.URL.Query().Get("end_date"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomById(roomId)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot get room from database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var reservation models.Reservation
	reservation.RoomId = roomId
	reservation.StartDate = startDate
	reservation.EndDate = endDate
	reservation.Room = room

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
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

type roomJsonResp struct {
	Rooms   []models.Room `json:"rooms"`
	Message string        `json:"message"`
}

func (m *Repository) GetRoomList(w http.ResponseWriter, r *http.Request) {
	rooms, err := m.DB.GetRooms()
	roomsResp := roomJsonResp{}
	status := http.StatusOK
	if err != nil {
		roomsResp.Message = err.Error()
		roomsResp.Rooms = []models.Room{}
		status = http.StatusInternalServerError
	}

	roomsResp.Rooms = rooms
	roomsResp.Message = "success"
	out, err := json.MarshalIndent(roomsResp, "", "  ")
	if err != nil {
		m.App.ErrorLog.Println("cannot marshal rooms to json")
		roomsResp.Message = err.Error()
		roomsResp.Rooms = []models.Room{}
		status = http.StatusInternalServerError
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(out)
}

// func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
// 	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
// }

// func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
// 	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
// }
