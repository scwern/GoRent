package car

import "github.com/google/uuid"

type Car struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Model       string    `json:"model" db:"model"`
	Brand       string    `json:"brand" db:"brand"`
	Year        int       `json:"year" db:"year"`
	PricePerDay float64   `json:"price_per_day" db:"price_per_day"`
	IsAvailable bool      `json:"is_available" db:"is_available"`
}
