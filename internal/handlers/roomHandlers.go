package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/helpers"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/models"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/render"
)

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

	render.Template(w, r, "roomDetails.page.tmpl", &models.TemplateData{
		Data: data,
	})
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
