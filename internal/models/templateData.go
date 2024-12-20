package models

import "github.com/thanhphuocnguyen/go-bookings-app/internal/forms"

type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Flash           string
	Warning         string
	Error           string
	CurrentUser     User
	Form            *forms.Form
	IsAuthenticated int
}
