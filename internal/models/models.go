package models

import (
	"time"
)

type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Phone       string
	Password    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	AccessLevel int
}

type Reservation struct {
	ID        int
	UserId    int
	RoomId    int
	FirstName string
	LastName  string
	Email     string
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Phone     string
	Room      Room
	User      User
}

type Restriction struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Room struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RoomRestriction struct {
	ID            int
	RoomId        int64
	RestrictionId int64
	ReservationId int64
	StartDate     time.Time
	EndDate       time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Room
	Restriction   Restriction
}