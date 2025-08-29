package handlers

import (
	"GoRent/internal/domain/user"
	"GoRent/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdminHandler struct {
	adminService service.AdminService
}

func NewAdminHandler(adminService service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// ChangeUserRole godoc
// @Summary      Сменить роль пользователя
// @Description  Доступно только admin
// @Tags         admin
// @Security     ApiKeyAuth
// @Accept       json
// @Param        id      path  string                  true  "ID пользователя"
// @Param        request body  user.ChangeRoleRequest  true  "Новая роль"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /admin/users/{id}/role [put]
func (h *AdminHandler) ChangeUserRole(c *gin.Context) {
	userID := c.Param("id")

	var req user.ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.adminService.ChangeUserRole(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListUsers godoc
// @Summary      Получить список пользователей
// @Description  Доступно только admin
// @Tags         admin
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {array}  object
// @Failure      403  {object}  map[string]string
// @Router       /admin/users [get]
func (h *AdminHandler) ListUsers(c *gin.Context) {
	users, err := h.adminService.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get users"})
		return
	}

	var response []gin.H
	for _, u := range users {
		response = append(response, gin.H{
			"id":    u.ID,
			"name":  u.Name,
			"email": u.Email,
			"role":  u.Role,
		})
	}

	c.JSON(200, response)
}
