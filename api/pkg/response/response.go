package response

import (
	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Success bool        `json:"success"`
	Data    T           `json:"data,omitempty"`
	Error   *ErrorItem  `json:"error,omitempty"`
	Meta    *MetaInfo   `json:"meta,omitempty"`
}

type ErrorItem struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type MetaInfo struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

func Success[T any](c *gin.Context, statusCode int, data T) {
	c.JSON(statusCode, Response[T]{
		Success: true,
		Data:    data,
	})
}

func SuccessWithMeta[T any](c *gin.Context, statusCode int, data T, meta MetaInfo) {
	c.JSON(statusCode, Response[T]{
		Success: true,
		Data:    data,
		Meta:    &meta,
	})
}

func Error(c *gin.Context, statusCode int, code string, message string) {
	c.JSON(statusCode, Response[any]{
		Success: false,
		Error: &ErrorItem{
			Code:    code,
			Message: message,
		},
	})
}
