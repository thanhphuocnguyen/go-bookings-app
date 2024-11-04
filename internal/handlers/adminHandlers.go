package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/forms"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/helpers"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/models"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/render"
)

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "adminDashboard.page.tmpl", &models.TemplateData{})
}

func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllNewReservations()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot get reservations from database")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}
	dataMap := make(map[string]interface{})
	dataMap["reservations"] = reservations

	render.Template(w, r, "adminNewReservations.page.tmpl", &models.TemplateData{
		Data: dataMap,
	})
}

func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllReservations()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot get reservations from database")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}
	dataMap := make(map[string]interface{})
	dataMap["reservations"] = reservations
	render.Template(w, r, "adminAllReservations.page.tmpl", &models.TemplateData{
		Data: dataMap,
	})
}

func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	next := now.AddDate(0, 1, 0)
	prev := now.AddDate(0, -1, 0)

	prevMonth, prevMonthYear, nextMonth, nextMonthYear := prev.Month(), prev.Year(), next.Month(), next.Year()
	strMap := make(map[string]string)

	strMap["next_month"] = strconv.Itoa(int(nextMonth))
	strMap["next_month_year"] = strconv.Itoa(nextMonthYear)
	strMap["prev_month"] = strconv.Itoa(int(prevMonth))
	strMap["prev_month_year"] = strconv.Itoa(prevMonthYear)
	curYear, curMonth, _ := now.Date()
	strMap["cur_month"] = strconv.Itoa(int(curMonth))
	strMap["cur_month_year"] = strconv.Itoa(curYear)
	firstOfMonth := time.Date(curYear, curMonth, 1, 0, 0, 0, 0, time.Now().Location())
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()
	dataMap := make(map[string]interface{})
	dataMap["now"] = now

	rooms, err := m.DB.GetRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	type roomRestrictionResult struct {
		roomID           int
		roomRestrictions []models.RoomRestriction
		err              error
	}

	results := make(chan roomRestrictionResult)
	for _, room := range rooms {
		go func(room models.Room) {
			roomRestrictions, err := m.DB.GetRoomRestrictionsForRoomByDate(room.ID, firstOfMonth, lastOfMonth)
			results <- roomRestrictionResult{roomID: room.ID, roomRestrictions: roomRestrictions, err: err}
		}(room)
	}

	for range rooms {
		result := <-results
		if result.err != nil {
			helpers.ServerError(w, result.err)
			return
		}

		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMonth; !d.After(lastOfMonth); d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-02")] = 0
			blockMap[d.Format("2006-01-02")] = 0
		}

		for _, y := range result.roomRestrictions {
			if y.ReservationId > 0 {
				for d := y.StartDate; !d.After(y.EndDate); d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = y.ReservationId
				}
			} else {
				blockMap[y.StartDate.Format("2006-01-2")] = y.ID
			}
		}

		dataMap[fmt.Sprintf("block_map_%d", result.roomID)] = blockMap
		dataMap[fmt.Sprintf("reservation_map_%d", result.roomID)] = reservationMap
		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", result.roomID), blockMap)
	}

	dataMap["rooms"] = rooms

	render.Template(w, r, "adminReservationsCalendar.page.tmpl", &models.TemplateData{
		StringMap: strMap,
		Data:      dataMap,
		IntMap:    intMap,
	})
}

func (m *Repository) AdminPostReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year, _ := strconv.Atoi(r.Form.Get("y"))
	month, _ := strconv.Atoi(r.Form.Get("m"))
	m.App.InfoLog.Println(year, month)
	form := forms.New(r.PostForm)

	rooms, err := m.DB.GetRooms()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var wg sync.WaitGroup

	for _, room := range rooms {
		blocks := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", room.ID)).(map[string]int)
		for k, v := range blocks {
			if v > 0 && !form.Has(fmt.Sprintf("remove_block_%d_%s", room.ID, k)) {
				wg.Add(1)
				go func() {
					err := m.DB.RemoveBlockById(v)
					defer wg.Done()
					if err != nil {
						helpers.ServerError(w, err)
						return
					}
				}()
			}
		}
	}

	wg.Wait()

	for k := range r.PostForm {
		if strings.HasPrefix(k, "add_block") {
			wg.Add(1)
			go func(k string) {
				defer wg.Done()
				exploded := strings.Split(k, "_")
				roomId, _ := strconv.Atoi(exploded[2])
				startDate, _ := time.Parse("2006-01-2", exploded[3])
				err := m.DB.InsertBlockForRoom(roomId, startDate)
				if err != nil {
					log.Println(err)
				}
			}(k)
		}
	}

	wg.Wait()

	m.App.Session.Put(r.Context(), "flash", "Save calendar successfully!")
	http.Redirect(w, r, "/admin/reservations-calendar", http.StatusSeeOther)
}

func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot parse reservation id")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	reservation, err := m.DB.GetReservationById(id)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot get reservation from database")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	strMap := make(map[string]string)
	strMap["from"] = r.URL.Query().Get("from")
	dataMap := make(map[string]interface{})
	dataMap["reservation"] = reservation

	render.Template(w, r, "adminShowReservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      dataMap,
		StringMap: strMap,
	})
}

func (m *Repository) AdminEditReservation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot parse reservation id")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	err = r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot parse form")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	reservation, err := m.DB.GetReservationById(id)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot get reservation from database")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	// Check if the form is valid
	f := forms.New(r.PostForm)
	f.Required("first_name", "last_name", "email", "phone")
	f.MinLength("first_name", 3, r)
	f.MinLength("last_name", 3, r)
	f.MinLength("phone", 10, r)
	f.IsEmail("email")

	dataMap := make(map[string]interface{})
	dataMap["reservation"] = reservation

	if !f.Valid() {
		m.App.Session.Put(r.Context(), "error", "Form is not valid")
		render.Template(w, r, "adminShowReservation.page.tmpl", &models.TemplateData{
			Form: f,
			Data: dataMap,
		})
		return
	}

	// Update the reservation
	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")
	// Update the reservation
	err = m.DB.UpdateReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot update reservation")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Reservation updated")
	from := r.URL.Query().Get("from")
	fmt.Println(from)
	if from == "all" {
		http.Redirect(w, r, "/admin/reservations-all", http.StatusSeeOther)
	} else if from == "new" {
		http.Redirect(w, r, "/admin/reservations-new", http.StatusSeeOther)
	} else {
		render.Template(w, r, "adminShowReservation.page.tmpl", &models.TemplateData{
			Form: f,
			Data: dataMap,
		})
	}
}

func (m *Repository) AdminProcessedReservation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	from := r.URL.Query().Get("from")

	redirectPath := "/admin/dashboard"
	if from == "all" {
		redirectPath = "/admin/reservations-all"
	} else if from == "new" {
		redirectPath = "/admin/reservations-new"
	}

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot parse reservation id")
		http.Redirect(w, r, redirectPath, http.StatusTemporaryRedirect)
		return
	}

	processedStr := r.URL.Query().Get("processed")

	if processedStr == "" {
		m.App.Session.Put(r.Context(), "error", "Cannot parse processed value")
		http.Redirect(w, r, redirectPath, http.StatusTemporaryRedirect)
		return
	}

	processed := processedStr == "true"
	err = m.DB.ProcessReservation(id, processed)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot update reservation")
		http.Redirect(w, r, redirectPath, http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Reservation processed")
	http.Redirect(w, r, redirectPath, http.StatusSeeOther)
}

func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot parse reservation id")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	err = m.DB.DeleteReservation(id)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot delete reservation")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Reservation deleted")
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}
