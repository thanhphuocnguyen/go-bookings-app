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
	UserTable            = "users"
)

func (m *pgRepository) AllUsers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`select * from %s`, UserTable)
	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		log.Println(err)
		return []models.User{}, err
	}

	users := []models.User{}

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.AccessLevel, &user.CreatedAt, &user.UpdatedAt)
		users = append(users, user)
		if err != nil {
			log.Println(err)
			return []models.User{}, err
		}
	}

	return users, nil
}

func (m *pgRepository) InsertReservation(res *models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	sql := fmt.Sprintf(`insert into %s 
		(user_id, room_id, email, first_name, last_name, phone, start_date, end_date) 
		values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`, ReservationTable)
	var newId int
	err := m.DB.QueryRowContext(ctx, sql, res.UserId, res.RoomId, res.Email, res.FirstName, res.LastName, res.Phone, res.StartDate, res.EndDate).Scan(&newId)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	log.Println("New reservation id: ", newId)

	return newId, nil
}

func (m *pgRepository) InsertRoomRestriction(res *models.RoomRestriction) error {
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

func (m *pgRepository) CheckIfRoomAvailableByDate(roomId int, start, end time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`
		select count(id) 
		from %s
			where room_id = $1 and $2 <= end_date and $3 >= start_date
	`, RoomRestrictionTable)
	var numRows int

	err := m.DB.QueryRowContext(ctx, query, roomId, start, end).Scan(&numRows)
	log.Println("numRows: ", numRows)
	if err != nil {
		log.Println(err)
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

func (m *pgRepository) SearchAvailabilityInRange(start, end time.Time) ([]models.Room, error) {
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

func (m *pgRepository) GetRoomById(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`select id, name, description, slug, price, created_at, updated_at from %s where id = $1`, RoomTable)
	var room models.Room
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&room.ID, &room.Name, &room.Description, &room.Slug, &room.Price, &room.CreatedAt, &room.UpdatedAt)

	if err != nil {
		return room, err
	}

	return room, nil
}

func (m *pgRepository) GetRoomBySlug(slug string) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := fmt.Sprintf(`select id, name, description, slug, price, created_at, updated_at from %s where slug = $1`, RoomTable)
	var room models.Room
	err := m.DB.QueryRowContext(ctx, query, slug).Scan(&room.ID, &room.Name, &room.Description, &room.Slug, &room.Price, &room.CreatedAt, &room.UpdatedAt)

	if err != nil {
		log.Println(err)
		return room, err
	}

	return room, nil
}

func (m *pgRepository) GetRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`select * from %s`, RoomTable)

	rooms := []models.Room{}
	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		log.Println(err)
		return []models.Room{}, err
	}
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.Name, &room.Description, &room.Slug, &room.Price, &room.CreatedAt, &room.UpdatedAt)
		if err != nil {
			return []models.Room{}, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}
