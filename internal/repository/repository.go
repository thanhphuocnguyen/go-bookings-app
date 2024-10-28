package repository

import (
	"time"

	"github.com/thanhphuocnguyen/go-bookings-app/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res *models.Reservation) (int64, error)
	InsertRoomRestriction(res *models.RoomRestriction) error

	CheckIfRoomAvailableByDate(roomId int, start, end time.Time) (bool, error)
	SearchAvailabilityInRange(start, end time.Time) ([]models.Room, error)

	GetRoomById(id int) (models.Room, error)
}
