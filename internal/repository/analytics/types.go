package analytics

type BrandRevenue struct {
	Brand   string  `db:"brand"`
	Revenue float64 `db:"revenue"`
	Count   int     `db:"rental_count"`
}

type RentalStats struct {
	TotalRentals    int     `db:"total_rentals"`
	TotalRevenue    float64 `db:"total_revenue"`
	AvgRentalDays   float64 `db:"avg_rental_days"`
	AvgDailyRevenue float64 `db:"avg_daily_revenue"`
}
