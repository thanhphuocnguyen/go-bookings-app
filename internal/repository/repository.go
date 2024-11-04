package repository

import (
	"time"

	"github.com/thanhphuocnguyen/go-bookings-app/internal/models"
)

type DatabaseRepo interface {
	//Reservations
	GetReservationById(id int) (models.Reservation, error)
	InsertReservation(res *models.Reservation) (int, error)
	InsertRoomRestriction(res *models.RoomRestriction) error
	CheckIfRoomAvailableByDate(roomId int, start, end time.Time) (bool, error)
	SearchAvailabilityInRange(start, end time.Time) ([]models.Room, error)
	GetRoomRestrictionsForRoomByDate(roomId int, start, end time.Time) ([]models.RoomRestriction, error)

	//Rooms
	GetRoomById(id int) (models.Room, error)
	GetRoomBySlug(slug string) (models.Room, error)
	GetRooms() ([]models.Room, error)

	//Users
	AllUsers() ([]models.User, error)
	GetUserById(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, error)

	//Admin
	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	ProcessReservation(id int, processed bool) error
	DeleteReservation(id int) error
	UpdateReservation(u models.Reservation) error
	InsertBlockForRoom(id int, startDate time.Time) error
	RemoveBlockById(id int) error
}
