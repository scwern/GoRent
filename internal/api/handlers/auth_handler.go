package handlers

import (
	"GoRent/internal/domain/user"
	"GoRent/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary      Регистрация пользователя
// @Description  Создает нового пользователя с ролью client
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body  user.RegisterRequest  true  "Данные регистрации"
// @Success      201      {object}  user.User
// @Failure      400      {object}  map[string]string
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req user.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

// Login godoc
// @Summary      Авторизация пользователя
// @Description  Вход в систему, возвращает JWT токен
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body  user.LoginRequest  true  "Данные входа"
// @Success      200      {object}  user.LoginResponse
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req user.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
