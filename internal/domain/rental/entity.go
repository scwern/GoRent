package rental

import (
	"github.com/google/uuid"
	"time"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusActive    Status = "active"
	StatusCompleted Status = "completed"
	StatusCanceled  Status = "canceled"
)

type Rental struct {
	ID         uuid.UUID `db:"id"`
	CarID      uuid.UUID `db:"car_id"`
	UserID     uuid.UUID `db:"user_id"`
	StartDate  time.Time `db:"start_date"`
	EndDate    time.Time `db:"end_date"`
	TotalPrice float64   `db:"total_price"`
	Status     Status    `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
}
