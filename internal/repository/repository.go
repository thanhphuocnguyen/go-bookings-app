package repository

import (
	"time"

	"github.com/thanhphuocnguyen/go-bookings-app/internal/models"
)

type DatabaseRepo interface {
	//Reservations
	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationById(id int) (models.Reservation, error)
	UpdateReservation(u models.Reservation) error
	DeleteReservation(id int) error
	ProcessReservation(id int, processed bool) error
	InsertReservation(res *models.Reservation) (int, error)
	InsertRoomRestriction(res *models.RoomRestriction) error
	CheckIfRoomAvailableByDate(roomId int, start, end time.Time) (bool, error)
	SearchAvailabilityInRange(start, end time.Time) ([]models.Room, error)

	//Rooms
	GetRoomById(id int) (models.Room, error)
	GetRoomBySlug(slug string) (models.Room, error)
	GetRooms() ([]models.Room, error)

	//Users
	AllUsers() ([]models.User, error)
	GetUserById(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, error)
}
