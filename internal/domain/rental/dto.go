package rental

import (
	"github.com/google/uuid"
	"time"
)

type CreateRequest struct {
	CarID     uuid.UUID `json:"car_id" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

type Response struct {
	ID         string    `json:"id"`
	CarID      string    `json:"car_id"`
	UserID     string    `json:"user_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	TotalPrice float64   `json:"total_price"`
	Status     Status    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type UpdateStatusRequest struct {
	Status Status `json:"status" binding:"required,oneof=pending active completed canceled"`
}
