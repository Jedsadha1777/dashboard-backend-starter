package handlers

import (
	"dashboard-starter/internal/application/dto"
	"dashboard-starter/internal/application/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	deviceService *services.DeviceApplicationService
}

func NewDeviceHandler(deviceService *services.DeviceApplicationService) *DeviceHandler {
	return &DeviceHandler{
		deviceService: deviceService,
	}
}

func (h *DeviceHandler) AuthenticateDevice(c *gin.Context) {
	var input dto.DeviceAuthInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
		})
		return
	}

	response, err := h.deviceService.AuthenticateDevice(input)
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

func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Admin ID not found",
		})
		return
	}

	var input struct {
		DeviceID string `json:"device_id" binding:"required"`
		Name     string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
		})
		return
	}

	device, apiKey, err := h.deviceService.RegisterDevice(input.DeviceID, input.Name)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "device ID already exists" {
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
		Data: gin.H{
			"id":         device.ID,
			"device_id":  device.DeviceID,
			"name":       device.Name,
			"api_key":    apiKey,
			"status":     device.Status,
			"created_by": adminID,
		},
	})
}

func (h *DeviceHandler) ListDevices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	devices, total, err := h.deviceService.ListDevices(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to retrieve devices: " + err.Error(),
		})
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    devices,
		Meta: gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}

func (h *DeviceHandler) GetDevice(c *gin.Context) {
	id := c.Param("id")
	deviceID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid device ID",
		})
		return
	}

	device, err := h.deviceService.GetDevice(uint(deviceID))
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Device not found",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    device,
	})
}

func (h *DeviceHandler) ResetAPIKey(c *gin.Context) {
	id := c.Param("id")
	deviceID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid device ID",
		})
		return
	}

	device, newAPIKey, err := h.deviceService.ResetDeviceAPIKey(uint(deviceID))
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Device not found",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: gin.H{
			"id":        device.ID,
			"device_id": device.DeviceID,
			"name":      device.Name,
			"api_key":   newAPIKey,
			"message":   "API key reset successfully",
		},
	})
}

func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: gin.H{
			"message":   "Update device functionality needs to be implemented",
			"device_id": id,
		},
	})
}

func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: gin.H{
			"message":   "Delete device functionality needs to be implemented",
			"device_id": id,
		},
	})
}
