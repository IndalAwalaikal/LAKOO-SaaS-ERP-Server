package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"lakoo/backend/internal/dto"
	"lakoo/backend/internal/service"
	"lakoo/backend/pkg/response"
)

type FinanceHandler struct {
	service service.FinanceService
}

func NewFinanceHandler(uu service.FinanceService) *FinanceHandler {
	return &FinanceHandler{service: uu}
}

func (h *FinanceHandler) Record(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	userID, _ := c.Get("user_id")

	var req dto.FinanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "BAD_REQUEST", "Invalid input format")
		return
	}

	res, err := h.service.RecordTransaction(tenantID.(string), userID.(string), &req)
	if err != nil {
		response.Error(c, 422, "UNPROCESSABLE_ENTITY", err.Error())
		return
	}

	response.Success(c, 201, res)
}

func (h *FinanceHandler) Delete(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")

	err := h.service.DeleteTransaction(id, tenantID.(string))
	if err != nil {
		response.Error(c, 404, "NOT_FOUND", "Transaction not found or unauthorized")
		return
	}
	response.Success(c, 200, gin.H{"message": "Transaction deleted"})
}

func (h *FinanceHandler) List(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")

	txs, err := h.service.GetTransactions(tenantID.(string))
	if err != nil {
		response.Error(c, 500, "INTERNAL_SERVER_ERROR", err.Error())
		return
	}
	
	meta := response.MetaInfo{ Total: len(txs), Page: 1, Limit: len(txs) }
	response.SuccessWithMeta(c, 200, txs, meta)
}

func (h *FinanceHandler) Summary(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")

	start := c.Query("start_date")
	end := c.Query("end_date")

	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	if start != "" {
		if parsed, err := time.Parse("2006-01-02", start); err == nil {
			startDate = parsed
		}
	}
	if end != "" {
		if parsed, err := time.Parse("2006-01-02", end); err == nil {
			// Include the entire end day
			endDate = parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
	}

	summary, err := h.service.GetSummary(tenantID.(string), startDate, endDate)
	if err != nil {
		response.Error(c, 500, "INTERNAL_SERVER_ERROR", err.Error())
		return
	}

	response.Success(c, 200, summary)
}
