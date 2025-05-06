package controllers

import (
	"dashboard-starter/models"
	"dashboard-starter/services"
	"dashboard-starter/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// ListUsers handles the request to list users with pagination and search
func ListUsers(c *gin.Context) {
	search := c.Query("search")

	// Get page and limit with defaults
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := services.GetUsers(search, page, limit)
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

// CreateUser handles the request to create a new user
func CreateUser(c *gin.Context) {
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

	user, err := services.CreateUser(&input)
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

	user, err := services.GetUserByID(id)
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

	user, err := services.UpdateUser(id, &input)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "record not found" {
			statusCode = http.StatusNotFound
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
	id := c.Param("id")

	err := services.DeleteUser(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "record not found" {
			statusCode = http.StatusNotFound
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
