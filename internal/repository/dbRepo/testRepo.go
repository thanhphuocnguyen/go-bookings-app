package dbRepo

import (
	"errors"
	"time"

	"github.com/thanhphuocnguyen/go-bookings-app/internal/models"
)

func (m *testDbRepo) AllUsers() ([]models.User, error) {
	return []models.User{}, nil
}

func (m *testDbRepo) InsertReservation(res *models.Reservation) (int, error) {
	if res.RoomId == 2 {
		return 0, errors.New("some error")
	}
	return 0, nil
}

func (m *testDbRepo) InsertRoomRestriction(res *models.RoomRestriction) error {
	if res.RoomId == 1000 {
		return errors.New("some error")
	}
	return nil
}

func (m *testDbRepo) CheckIfRoomAvailableByDate(roomId int, start, end time.Time) (bool, error) {
	if roomId == 2 {
		return false, errors.New("some error")
	}

	return true, nil
}

func (m *testDbRepo) SearchAvailabilityInRange(start, end time.Time) ([]models.Room, error) {
	return []models.Room{}, nil
}

func (m *testDbRepo) GetRoomById(id int) (models.Room, error) {
	if id >= 2 {
		return models.Room{}, errors.New("some error")
	}
	return models.Room{}, nil
}

func (m *testDbRepo) GetRoomBySlug(slug string) (models.Room, error) {
	return models.Room{}, nil
}

func (m *testDbRepo) GetRooms() ([]models.Room, error) {
	return []models.Room{}, nil
}
