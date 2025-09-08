package handlers

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

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

func ErrorResponse(errorMsg string) Response {
	return Response{
		Success: false,
		Error:   errorMsg,
	}
}

func RespondWithError(c *gin.Context, statusCode int, errorMsg string) {
	c.JSON(statusCode, ErrorResponse(errorMsg))
}

func RespondWithSuccess(c *gin.Context, statusCode int, data interface{}, meta ...interface{}) {
	c.JSON(statusCode, SuccessResponse(data, meta...))
}
