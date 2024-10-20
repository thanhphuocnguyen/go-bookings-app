package handlers

import (
	"net/http"

	"github.com/thanhphuocnguyen/go-bookings-app/pkg/config"
	"github.com/thanhphuocnguyen/go-bookings-app/pkg/models"
	"github.com/thanhphuocnguyen/go-bookings-app/pkg/render"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

func InitRepo(r *Repository) {
	Repo = r
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func (m *Repository) NewTest(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "test.page.tmpl", &models.TemplateData{})
}
