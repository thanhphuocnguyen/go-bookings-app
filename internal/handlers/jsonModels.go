package handlers

import "github.com/thanhphuocnguyen/go-bookings-app/internal/models"

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
type roomJsonResp struct {
	Rooms   []models.Room `json:"rooms"`
	Message string        `json:"message"`
}
