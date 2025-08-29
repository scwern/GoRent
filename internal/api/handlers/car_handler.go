package handlers

import (
	"GoRent/internal/domain/car"
	"GoRent/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type CarHandler struct {
	carService service.CarService
}

func NewCarHandler(carService service.CarService) *CarHandler {
	return &CarHandler{carService: carService}
}

// CreateCar godoc
// @Summary      Создать автомобиль
// @Description  Доступно только manager и admin
// @Tags         cars
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        request  body  car.CreateRequest  true  "Данные автомобиля"
// @Success      201      {object}  car.Response
// @Failure      400      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Router       /cars [post]
func (h *CarHandler) CreateCar(c *gin.Context) {
	var input car.CreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input data"})
		return
	}

	createdCar, err := h.carService.CreateCar(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdCar)
}

// GetCar godoc
// @Summary      Получить автомобиль по ID
// @Tags         cars
// @Produce      json
// @Param        id  path  string  true  "ID автомобиля"
// @Success      200  {object}  car.Response
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /cars/{id} [get]
func (h *CarHandler) GetCar(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid car ID"})
		return
	}

	carItem, err := h.carService.GetCar(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
		return
	}

	c.JSON(http.StatusOK, carItem)
}

// GetAllCars godoc
// @Summary      Получить список автомобилей
// @Description  Возвращает список автомобилей с возможностью фильтрации
// @Tags         cars
// @Produce      json
// @Param        available   query  bool    false  "Фильтр по доступности"
// @Param        brand       query  string  false  "Фильтр по бренду"
// @Param        min_price   query  number  false  "Минимальная цена"
// @Param        max_price   query  number  false  "Максимальная цена"
// @Success      200         {array}  car.Response
// @Router       /cars [get]
func (h *CarHandler) GetAllCars(c *gin.Context) {
	filters := make(map[string]interface{})

	if available := c.Query("available"); available != "" {
		filters["available"] = available == "true"
	}

	if brand := c.Query("brand"); brand != "" {
		filters["brand"] = brand
	}

	if minPrice := c.Query("min_price"); minPrice != "" {
		if price, err := parseFloat(minPrice); err == nil {
			filters["min_price"] = price
		}
	}

	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if price, err := parseFloat(maxPrice); err == nil {
			filters["max_price"] = price
		}
	}

	cars, err := h.carService.GetAllCars(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cars)
}

// UpdateCar godoc
// @Summary      Обновить автомобиль
// @Description  Доступно только manager и admin
// @Tags         cars
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id       path  string            true  "ID автомобиля"
// @Param        request  body  car.UpdateRequest true  "Данные для обновления"
// @Success      200      {object}  car.Response
// @Failure      400      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Router       /cars/{id} [put]
func (h *CarHandler) UpdateCar(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid car ID"})
		return
	}

	var input car.UpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input data"})
		return
	}

	updatedCar, err := h.carService.UpdateCar(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCar)
}

// DeleteCar godoc
// @Summary      Удалить автомобиль
// @Description  Доступно только manager и admin
// @Tags         cars
// @Security     ApiKeyAuth
// @Param        id  path  string  true  "ID автомобиля"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /cars/{id} [delete]
func (h *CarHandler) DeleteCar(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid car ID"})
		return
	}

	if err := h.carService.DeleteCar(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}
