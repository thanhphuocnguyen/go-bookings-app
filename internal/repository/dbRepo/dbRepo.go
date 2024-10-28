package dbRepo

import (
	"database/sql"

	"github.com/thanhphuocnguyen/go-bookings-app/internal/config"
)

type PGRepository struct {
	App *config.AppConfig
	DB  *sql.DB
}

func InitPGRepository(app *config.AppConfig, db *sql.DB) *PGRepository {
	return &PGRepository{
		App: app,
		DB:  db,
	}
}
