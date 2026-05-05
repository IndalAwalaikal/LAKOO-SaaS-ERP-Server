package handler

import (
	"lakoo/backend/internal/service"
	"lakoo/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service service.NotificationService
}

func NewNotificationHandler(u service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: u}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	res, err := h.service.GetNotifications(tenantID.(string))
	if err != nil {
		response.Error(c, 500, "INTERNAL_SERVER_ERROR", err.Error())
		return
	}
	response.Success(c, 200, res)
}
