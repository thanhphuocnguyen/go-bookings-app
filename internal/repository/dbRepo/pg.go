package dbRepo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/thanhphuocnguyen/go-bookings-app/internal/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	ReservationTable     = "reservations"
	RoomRestrictionTable = "room_restrictions"
	RoomTable            = "rooms"
	UserTable            = "users"
)

// User services
func (m *pgRepository) AllUsers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`select id, first_name, last_name, email, password, access_level, created_at, updated_at from %s`, UserTable)
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
			log.Println("AllUsers", err)
			return []models.User{}, err
		}
	}

	return users, nil
}

func (m *pgRepository) GetUserById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user models.User
	query := fmt.Sprintf(`select id, email, phone, first_name, last_name, password, access_level from %s where id=$1`, UserTable)

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.Phone, &user.FirstName, &user.LastName, &user.AccessLevel)

	if err != nil {
		log.Println("GetUserById", err)
		return user, err
	}

	return user, nil
}

func (m *pgRepository) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`update %s set email=$1, phone=$2, first_name=$3, last_name=$4, password=$5, access_level=$6, updated_at=$7 where id=$8`, UserTable)
	_, err := m.DB.ExecContext(ctx, query, u.Email, u.Phone, u.FirstName, u.LastName, u.Password, u.AccessLevel, time.Now(), u.ID)

	if err != nil {
		log.Println("UpdateUser", err)
		return err
	}

	return nil
}

func (m *pgRepository) Authenticate(email, testPassword string) (int, string, error) {
	var id int
	var hashedPassword string
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf("select id, password from %s where email=$1", UserTable)
	err := m.DB.QueryRowContext(ctx, query, email).Scan(&id, &hashedPassword)

	if err != nil {
		log.Println("Authenticate", err)
		return 0, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return 0, "", errors.New("password does not match")
	} else if err != nil {
		return 0, "", err
	}
	return id, hashedPassword, nil
}

// Reservation actions
func (m *pgRepository) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`
		select 
			rs.id, rs.user_id, rs.room_id, rs.email, rs.first_name, rs.last_name, rs.phone, 
			rs.start_date, rs.end_date, rs.processed, rs.created_at, rs.updated_at, r.id, r.name, r.price
		from %s rs
		left join %s r on rs.room_id = r.id
	`, ReservationTable, RoomTable)

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		log.Println("AllReservations", err)
		return []models.Reservation{}, err
	}

	var reservations []models.Reservation

	for rows.Next() {
		var res models.Reservation
		err := rows.Scan(&res.ID, &res.UserId, &res.RoomId, &res.Email, &res.FirstName, &res.LastName, &res.Phone, &res.StartDate, &res.EndDate, &res.Processed, &res.CreatedAt, &res.UpdatedAt, &res.Room.ID, &res.Room.Name, &res.Room.Price)
		if err != nil {
			log.Println("AllReservations", err)
			return []models.Reservation{}, err
		}

		reservations = append(reservations, res)
	}

	return reservations, nil
}

func (m *pgRepository) AllNewReservations() ([]models.Reservation, error) {
	cxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`
		select 
			rs.id, rs.user_id, rs.room_id, rs.email, rs.first_name, rs.last_name, rs.phone, 
			rs.start_date, rs.end_date, rs.created_at, rs.updated_at, r.id, r.name, r.price
		from %s rs
		left join %s r on rs.room_id = r.id
		where rs.processed = false and rs.start_date > $1 and rs.end_date > $2
	`, ReservationTable, RoomTable)

	rows, err := m.DB.QueryContext(cxt, query, time.Now(), time.Now())

	if err != nil || rows.Err() != nil {
		log.Println("AllReservations", err)
		return []models.Reservation{}, err
	}

	var reservations []models.Reservation

	for rows.Next() {
		var res models.Reservation
		err := rows.Scan(&res.ID, &res.UserId, &res.RoomId, &res.Email, &res.FirstName, &res.LastName, &res.Phone, &res.StartDate, &res.EndDate, &res.CreatedAt, &res.UpdatedAt, &res.Room.ID, &res.Room.Name, &res.Room.Price)
		if err != nil {
			log.Println("AllReservations", err)
			return []models.Reservation{}, err
		}

		reservations = append(reservations, res)
	}

	return reservations, nil

}

func (m *pgRepository) GetReservationById(id int) (models.Reservation, error) {
	cxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`
		select 
			rs.id, rs.user_id, rs.room_id, rs.email, rs.first_name, rs.last_name, rs.phone, 
			rs.start_date, rs.end_date, rs.processed, rs.created_at, rs.updated_at, r.id, r.name, r.price
		from %s rs
		left join %s r on rs.room_id = r.id
		where rs.id = $1
	`, ReservationTable, RoomTable)
	var res models.Reservation
	err := m.DB.QueryRowContext(cxt, query, id).Scan(&res.ID, &res.UserId, &res.RoomId, &res.Email, &res.FirstName, &res.LastName, &res.Phone, &res.StartDate, &res.EndDate, &res.Processed, &res.CreatedAt, &res.UpdatedAt, &res.Room.ID, &res.Room.Name, &res.Room.Price)

	if err != nil {
		log.Println("GetReservationById", err)
		return res, err
	}

	return res, nil
}

func (m *pgRepository) UpdateReservation(u models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`update %s set first_name=$1, last_name=$2, email=$3, phone=$4, updated_at=$5 where id=$6`, ReservationTable)
	_, err := m.DB.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.Phone, time.Now(), u.ID)

	if err != nil {
		log.Println("UpdateReservation", err)
		return err
	}

	return nil
}

func (m *pgRepository) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`delete from %s where id=$1`, ReservationTable)
	_, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		log.Println("DeleteReservation", err)
		return err
	}

	return nil
}

func (m *pgRepository) ProcessReservation(id int, processed bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`update %s set processed=$1 where id=$2`, ReservationTable)
	_, err := m.DB.ExecContext(ctx, query, processed, id)

	if err != nil {
		log.Println("ProcessReservation", err)
		return err
	}

	return nil
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
		log.Println("InsertReservation", err)
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
		log.Println("InsertRoomRestriction", err)
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
			where room_id = $1 and $2 < end_date and $3 > start_date
	`, RoomRestrictionTable)
	var numRows int

	err := m.DB.QueryRowContext(ctx, query, roomId, start, end).Scan(&numRows)
	log.Println("numRows: ", numRows)
	if err != nil {
		log.Println("CheckIfRoomAvailableByDate", err)
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
		log.Println("SearchAvailabilityInRange", err)
		log.Println(err)
		return nil, err
	}

	var rooms []models.Room

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.Name)
		if err != nil {
			log.Println("SearchAvailabilityInRange", err)
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
		log.Println("GetRoomById", err)
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
		log.Println("GetRoomBySlug", err)
		log.Println(err)
		return room, err
	}

	return room, nil
}

func (m *pgRepository) GetRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`select id, name, description, slug, price, created_at, updated_at from %s`, RoomTable)

	rooms := []models.Room{}
	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		log.Println("GetRooms", err)
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
