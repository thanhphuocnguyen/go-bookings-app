package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/handlers"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/contact", handlers.Repo.Contact)
	// mux.Get("/generals-quarters", handlers.Repo.Generals)
	// mux.Get("/majors-suite", handlers.Repo.Majors)
	mux.Get("/rooms/{id}", handlers.Repo.Room)

	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.CreateReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.SearchAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	mux.Get("/book-room", handlers.Repo.BookRoom)
	mux.Get("/get-all-rooms", handlers.Repo.GetRoomList)

	// Users
	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Post("/user/login", handlers.Repo.PostLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)
	mux.Get("/user/register", handlers.Repo.ShowRegistration)
	mux.Post("/user/register", handlers.Repo.PostRegistration)

	// Admin
	mux.Route("/admin", func(r chi.Router) {
		r.Use(Auth)
		r.Get("/dashboard", handlers.Repo.AdminDashboard)
	})
	return mux
}
