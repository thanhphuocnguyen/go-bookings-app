package handlers

import (
	"fmt"
	"net/http"

	"github.com/thanhphuocnguyen/go-bookings-app/internal/forms"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/models"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/render"
)

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

	render.Template(w, r, "makeReservation.page.tmpl", &models.TemplateData{
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
		render.Template(w, r, "makeReservation.page.tmpl", &models.TemplateData{
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
	render.Template(w, r, "reservationSummary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: strMap,
	})
}
