package dbRepo

import (
	"database/sql"

	"github.com/thanhphuocnguyen/go-bookings-app/internal/config"
)

type pgRepository struct {
	App *config.AppConfig
	DB  *sql.DB
}

type testDbRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func InitPGRepository(app *config.AppConfig, db *sql.DB) *pgRepository {
	return &pgRepository{
		App: app,
		DB:  db,
	}
}

func InitTestingRepository(app *config.AppConfig, db *sql.DB) *testDbRepo {
	return &testDbRepo{
		App: app,
		DB:  db,
	}
}
