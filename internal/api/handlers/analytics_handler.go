package handlers

import (
	"GoRent/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type AnalyticsHandler struct {
	analyticsService service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

// GetProfit godoc
// @Summary      Получить прибыль
// @Description  Возвращает суммарную прибыль за период
// @Tags         analytics
// @Security     ApiKeyAuth
// @Produce      json
// @Param        from  query  string  true  "Начальная дата (RFC3339)"
// @Param        to    query  string  true  "Конечная дата (RFC3339)"
// @Success      200   {object}  object
// @Failure      400   {object}  map[string]string
// @Failure      403   {object}  map[string]string
// @Router       /analytics/profit [get]
func (h *AnalyticsHandler) GetProfit(c *gin.Context) {
	fromDate, toDate, err := parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profit, err := h.analyticsService.GetProfit(c.Request.Context(), fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profit":    profit,
		"from_date": fromDate.Format(time.RFC3339),
		"to_date":   toDate.Format(time.RFC3339),
		"currency":  "USD",
	})
}

// GetPopularBrands godoc
// @Summary      Получить популярные бренды
// @Description  Возвращает самые популярные бренды по доходу
// @Tags         analytics
// @Security     ApiKeyAuth
// @Produce      json
// @Param        from   query  string  true   "Начальная дата (RFC3339)"
// @Param        to     query  string  true   "Конечная дата (RFC3339)"
// @Param        limit  query  int     false  "Лимит результатов"
// @Success      200    {object}  object
// @Failure      400    {object}  map[string]string
// @Failure      403    {object}  map[string]string
// @Router       /analytics/popular-brands [get]
func (h *AnalyticsHandler) GetPopularBrands(c *gin.Context) {
	fromDate, toDate, err := parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if limitInt, err := strconv.Atoi(limitStr); err == nil && limitInt > 0 {
			limit = limitInt
		}
	}

	brands, err := h.analyticsService.GetPopularBrands(c.Request.Context(), fromDate, toDate, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"brands":    brands,
		"from_date": fromDate.Format(time.RFC3339),
		"to_date":   toDate.Format(time.RFC3339),
		"limit":     limit,
	})
}

// GetRentalStats godoc
// @Summary      Получить статистику аренд
// @Description  Возвращает общую статистику по арендам
// @Tags         analytics
// @Security     ApiKeyAuth
// @Produce      json
// @Param        from  query  string  true  "Начальная дата (RFC3339)"
// @Param        to    query  string  true  "Конечная дата (RFC3339)"
// @Success      200   {object}  object
// @Failure      400   {object}  map[string]string
// @Failure      403   {object}  map[string]string
// @Router       /analytics/stats [get]
func (h *AnalyticsHandler) GetRentalStats(c *gin.Context) {
	fromDate, toDate, err := parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stats, err := h.analyticsService.GetRentalStats(c.Request.Context(), fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats":     stats,
		"from_date": fromDate.Format(time.RFC3339),
		"to_date":   toDate.Format(time.RFC3339),
	})
}

func parseDateRange(c *gin.Context) (time.Time, time.Time, error) {
	fromStr := c.Query("from")
	toStr := c.Query("to")

	if fromStr == "" && toStr == "" {
		now := time.Now()
		fromDate := time.Date(2025, now.Month(), 1, 0, 0, 0, 0, now.Location())
		toDate := fromDate.AddDate(0, 1, -1)
		return fromDate, toDate, nil
	}

	if fromStr == "" || toStr == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("both from and to parameters are required when one is specified")
	}

	fromDate, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid from date format. Use RFC3339 (e.g., 2025-01-01T00:00:00Z)")
	}

	toDate, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid to date format. Use RFC3339 (e.g., 2025-12-31T23:59:59Z)")
	}

	if fromDate.After(toDate) {
		return time.Time{}, time.Time{}, fmt.Errorf("from date cannot be after to date")
	}

	return fromDate, toDate, nil
}
