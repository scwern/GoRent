package car

type CreateRequest struct {
	Model       string  `json:"model" binding:"required,min=1,max=100"`
	Brand       string  `json:"brand" binding:"required,min=1,max=50"`
	Year        int     `json:"year" binding:"required,gte=1900,lte=2025"`
	PricePerDay float64 `json:"price_per_day" binding:"required,gt=0"`
}

type UpdateRequest struct {
	Model       *string  `json:"model" binding:"omitempty,min=1,max=100"`
	Brand       *string  `json:"brand" binding:"omitempty,min=1,max=50"`
	Year        *int     `json:"year" binding:"omitempty,gte=1900,lte=2025"`
	PricePerDay *float64 `json:"price_per_day" binding:"omitempty,gt=0"`
	IsAvailable *bool    `json:"is_available"`
}

type Response struct {
	ID          string  `json:"id"`
	Model       string  `json:"model"`
	Brand       string  `json:"brand"`
	Year        int     `json:"year"`
	PricePerDay float64 `json:"price_per_day"`
	IsAvailable bool    `json:"is_available"`
}
