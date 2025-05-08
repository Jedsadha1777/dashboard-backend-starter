package controllers

import (
	"dashboard-starter/models"
	"dashboard-starter/services"
	"dashboard-starter/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListUsers handles the request to list users with pagination and search
func ListUsers(c *gin.Context) {
	// Convert query parameters to PaginationParams
	var params utils.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		// Use default values
		params = utils.NewPaginationParams()
	}

	// Call service
	userService := services.NewUserService()
	users, pagination, err := userService.GetUsers(params)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to retrieve users: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users,
		Meta:    pagination,
	})
}

// CreateUser handles the request to create a new user
func CreateUser(c *gin.Context) {
	// Get admin ID from context
	adminID, _ := c.Get("admin_id")

	var input models.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
		})
		return
	}

	// Validate input
	if err := utils.ValidateStruct(input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Create service and call
	userService := services.NewUserService()
	user, err := userService.CreateUser(&input, adminID.(uint))

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to create user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    user,
	})
}

// GetUser handles the request to get a user by ID
func GetUser(c *gin.Context) {
	id := c.Param("id")

	// Create service and call
	userService := services.NewUserService()
	user, err := userService.GetByID(id)

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

// UpdateUser handles the request to update a user
func UpdateUser(c *gin.Context) {
	// Get admin ID from context
	adminID, _ := c.Get("admin_id")

	id := c.Param("id")

	var input models.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
		})
		return
	}

	// Validate input
	if err := utils.ValidateStruct(input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Create service and call
	userService := services.NewUserService()
	user, err := userService.UpdateUser(id, &input, adminID.(uint))

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "record not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "you don't have permission to update this user" {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, Response{
			Success: false,
			Error:   "Failed to update user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
	})
}

// DeleteUser handles the request to delete a user
func DeleteUser(c *gin.Context) {
	// Get admin ID from context
	adminID, _ := c.Get("admin_id")

	id := c.Param("id")

	// Create service and call
	userService := services.NewUserService()
	err := userService.DeleteUser(id, adminID.(uint))

	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "record not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "you don't have permission to delete this user" {
			statusCode = http.StatusForbidden
		}

		c.JSON(statusCode, Response{
			Success: false,
			Error:   "Failed to delete user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    gin.H{"message": "User deleted successfully"},
	})
}
