package repository

import (
	"time"

	"github.com/thanhphuocnguyen/go-bookings-app/internal/models"
)

type DatabaseRepo interface {
	InsertReservation(res *models.Reservation) (int, error)
	InsertRoomRestriction(res *models.RoomRestriction) error

	CheckIfRoomAvailableByDate(roomId int, start, end time.Time) (bool, error)
	SearchAvailabilityInRange(start, end time.Time) ([]models.Room, error)

	GetRoomById(id int) (models.Room, error)
	GetRoomBySlug(slug string) (models.Room, error)
	GetRooms() ([]models.Room, error)

	//Users
	AllUsers() ([]models.User, error)
	GetUserById(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, error)
}
