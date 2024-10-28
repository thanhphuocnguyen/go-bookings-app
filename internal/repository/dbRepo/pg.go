package dbRepo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/thanhphuocnguyen/go-bookings-app/internal/models"
)

const (
	ReservationTable     = "reservations"
	RoomRestrictionTable = "room_restrictions"
	RoomTable            = "rooms"
)

func (m *PGRepository) AllUsers() bool {
	return true
}

func (m *PGRepository) InsertReservation(res *models.Reservation) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	sql := fmt.Sprintf(`insert into %s 
		(id, user_id, room_id, start_date, end_date, phone) 
		values ($1, $2, $3, $4, $5, $6)`, ReservationTable)
	rs, err := m.DB.ExecContext(ctx, sql, 1, res.UserId, res.RoomId, res.StartDate, res.EndDate, res.Phone)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	newId, err := rs.LastInsertId()

	if err != nil {
		log.Println(err)
		return 0, err
	}

	log.Println(rs)
	return newId, nil
}

func (m *PGRepository) InsertRoomRestriction(res *models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	sql := fmt.Sprintf(`insert into %s 
		(room_id, restriction_id, reservation_id, start_date, end_date) 
		values ($1, $2, $3, $4, $5)`, RoomRestrictionTable)
	_, err := m.DB.ExecContext(ctx, sql, res.RoomId, res.RestrictionId, res.ReservationId, res.StartDate, res.EndDate)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *PGRepository) CheckIfRoomAvailableByDate(roomId int, start, end time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`
		select count(id) from %s
			where room_id = $1 and $2 < end_date and $3 > start_date
	`, RoomRestrictionTable)
	var numRows int
	err := m.DB.QueryRowContext(ctx, query, roomId, start, end).Scan(&numRows)

	if err != nil {
		log.Println(err)
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

func (m *PGRepository) SearchAvailabilityInRange(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`
		select id, name from %s 
			where id not in (select room_id from %s where $1 < end_date and $2 > start_date)
	`, RoomTable, RoomRestrictionTable)

	rows, err := m.DB.QueryContext(ctx, query, start, end)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var rooms []models.Room

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.Name)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}
