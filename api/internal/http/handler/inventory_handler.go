package handler

import (
	"github.com/gin-gonic/gin"
	"lakoo/backend/internal/dto"
	"lakoo/backend/internal/service"
	"lakoo/backend/pkg/response"
)

type InventoryHandler struct {
	service service.InventoryService
}

func NewInventoryHandler(uu service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: uu}
}

func (h *InventoryHandler) Adjust(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	userID, _ := c.Get("user_id")

	var req dto.InventoryMutationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "BAD_REQUEST", "Invalid input format")
		return
	}

	res, err := h.service.AdjustStock(tenantID.(string), userID.(string), &req)
	if err != nil {
		response.Error(c, 422, "UNPROCESSABLE_ENTITY", err.Error())
		return
	}

	response.Success(c, 201, res)
}

func (h *InventoryHandler) History(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	productID := c.Param("productId")

	history, err := h.service.GetProductHistory(tenantID.(string), productID)
	if err != nil {
		response.Error(c, 500, "INTERNAL_SERVER_ERROR", err.Error())
		return
	}
	
	meta := response.MetaInfo{ Total: len(history), Page: 1, Limit: len(history) }
	response.SuccessWithMeta(c, 200, history, meta)
}
