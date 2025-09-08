package handlers

import (
	"dashboard-starter/internal/application/dto"
	"dashboard-starter/internal/application/services"
	"dashboard-starter/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserHandler struct {
	userService *services.UserApplicationService
}

func NewUserHandler(userService *services.UserApplicationService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var input dto.UserRegistrationInput

	if err := c.ShouldBindBodyWith(&input, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
		})
		return
	}

	if err := utils.ValidateStruct(input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	response, err := h.userService.RegisterUser(input)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "email already exists" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    response,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var input dto.UserLoginInput

	if err := c.ShouldBindBodyWith(&input, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
		})
		return
	}

	if err := utils.ValidateStruct(input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	response, err := h.userService.LoginUser(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    response,
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "User ID not found",
		})
		return
	}

	user, err := h.userService.GetUserProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "User ID not found",
		})
		return
	}

	var input dto.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
		})
		return
	}

	if err := utils.ValidateStruct(input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	user, err := h.userService.UpdateUserProfile(userID.(uint), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to update profile: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
	})
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "User ID not found",
		})
		return
	}

	var input dto.UserChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
		})
		return
	}

	if err := utils.ValidateStruct(input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := h.userService.ChangePassword(userID.(uint), input); err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "current password is incorrect" {
			statusCode = http.StatusUnauthorized
		}
		c.JSON(statusCode, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    gin.H{"message": "Password updated successfully"},
	})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := h.userService.ListUsers(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to retrieve users: " + err.Error(),
		})
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users,
		Meta: gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Admin ID not found",
		})
		return
	}

	var input dto.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
		})
		return
	}

	if err := utils.ValidateStruct(input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	user, tempPassword, err := h.userService.CreateUserByAdmin(input, adminID.(uint))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "email already exists" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, Response{
			Success: false,
			Error:   "Failed to create user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data: gin.H{
			"user":               user,
			"temporary_password": tempPassword,
			"message":            "User created successfully. Please inform the user to change their password after first login.",
		},
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	user, err := h.userService.GetUserProfile(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Admin ID not found",
		})
		return
	}

	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	var input dto.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
		})
		return
	}

	// TODO: Implement UpdateUserByAdmin in user service
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: gin.H{
			"message":  "Update user functionality needs to be implemented",
			"user_id":  userID,
			"admin_id": adminID,
		},
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: gin.H{
			"message": "Delete user functionality needs to be implemented",
			"user_id": id,
		},
	})
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: gin.H{
			"message": "Reset password functionality needs to be implemented",
			"user_id": id,
		},
	})
}
