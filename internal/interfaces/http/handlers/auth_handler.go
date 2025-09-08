package handlers

import (
	"dashboard-starter/internal/application/dto"
	"dashboard-starter/internal/application/services"
	"dashboard-starter/internal/domain/shared/errors"
	"dashboard-starter/internal/interfaces/http/middleware"
	"dashboard-starter/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// @Summary Admin login
// @Description Authenticate admin user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body dto.LoginInput true "Login credentials"
// @Success 200 {object} Response{data=dto.LoginResponse} "Login successful"
// @Failure 400 {object} Response "Invalid input"
// @Failure 401 {object} Response "Invalid credentials"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input dto.LoginInput

	if err := c.ShouldBindBodyWith(&input, binding.JSON); err != nil {
		middleware.HandleError(c, errors.ValidationError(map[string]string{
			"input": err.Error(),
		}))
		return
	}

	if err := utils.ValidateStruct(input); err != nil {
		middleware.HandleError(c, errors.ValidationError(map[string]string{
			"validation": err.Error(),
		}))
		return
	}

	response, err := h.authService.LoginAdmin(input)
	if err != nil {
		middleware.HandleError(c, errors.InvalidCredentials())
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    response,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		middleware.HandleError(c, errors.ErrUnauthorized)
		return
	}

	if err := h.authService.LogoutAdmin(adminID.(uint)); err != nil {
		middleware.HandleError(c, errors.WrapError(err, errors.ErrCodeInternal, 500))
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    gin.H{"message": "Logged out successfully"},
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input dto.RefreshTokenInput

	if err := c.ShouldBindJSON(&input); err != nil {
		middleware.HandleError(c, errors.ValidationError(map[string]string{
			"input": err.Error(),
		}))
		return
	}

	response, err := h.authService.RefreshToken(input)
	if err != nil {
		middleware.HandleError(c, errors.TokenExpired())
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    response,
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		middleware.HandleError(c, errors.ErrUnauthorized)
		return
	}

	admin, err := h.authService.GetAdminProfile(adminID.(uint))
	if err != nil {
		middleware.HandleError(c, errors.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    admin,
	})
}
