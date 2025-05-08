package controllers

import "github.com/gin-gonic/gin"

// Response เป็นโครงสร้างมาตรฐานสำหรับการตอบกลับ API ทั้งหมด
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// NewSuccessResponse สร้าง Response สำหรับการทำงานที่สำเร็จ
func SuccessResponse(data interface{}, meta ...interface{}) Response {
	response := Response{
		Success: true,
		Data:    data,
	}

	if len(meta) > 0 && meta[0] != nil {
		response.Meta = meta[0]
	}

	return response
}

// NewErrorResponse สร้าง Response สำหรับการทำงานที่มีข้อผิดพลาด
func ErrorResponse(errorMsg string) Response {
	return Response{
		Success: false,
		Error:   errorMsg,
	}
}

// RespondWithError ส่งค่า error response กลับไปยัง client
func RespondWithError(c *gin.Context, statusCode int, errorMsg string) {
	c.JSON(statusCode, ErrorResponse(errorMsg))
}

// RespondWithSuccess ส่งค่า success response กลับไปยัง client
func RespondWithSuccess(c *gin.Context, statusCode int, data interface{}, meta ...interface{}) {
	c.JSON(statusCode, SuccessResponse(data, meta...))
}
