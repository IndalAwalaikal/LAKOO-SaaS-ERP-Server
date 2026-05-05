package handler

import (
	"github.com/gin-gonic/gin"
	"lakoo/backend/internal/dto"
	"lakoo/backend/internal/service"
	"lakoo/backend/pkg/response"
)

type CustomerHandler struct {
	service service.CustomerService
}

func NewCustomerHandler(uu service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: uu}
}

func (h *CustomerHandler) Create(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.Error(c, 403, "FORBIDDEN", "Tenant Context Missing")
		return
	}

	var req dto.CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "BAD_REQUEST", "Invalid input format")
		return
	}

	res, err := h.service.CreateCustomer(tenantID.(string), &req)
	if err != nil {
		response.Error(c, 422, "UNPROCESSABLE_ENTITY", err.Error())
		return
	}

	response.Success(c, 201, res)
}

func (h *CustomerHandler) Update(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")

	var req dto.CustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "BAD_REQUEST", "Invalid input format")
		return
	}

	res, err := h.service.UpdateCustomer(id, tenantID.(string), &req)
	if err != nil {
		if err.Error() == "customer not found" {
			response.Error(c, 404, "NOT_FOUND", "Customer not found")
		} else {
			response.Error(c, 400, "BAD_REQUEST", err.Error())
		}
		return
	}
	response.Success(c, 200, res)
}

func (h *CustomerHandler) Delete(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")

	err := h.service.DeleteCustomer(id, tenantID.(string))
	if err != nil {
		response.Error(c, 404, "NOT_FOUND", "Customer not found or unauthorized")
		return
	}
	response.Success(c, 200, gin.H{"message": "Customer successfully deleted"})
}

func (h *CustomerHandler) List(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")

	customers, err := h.service.GetCustomers(tenantID.(string))
	if err != nil {
		response.Error(c, 500, "INTERNAL_SERVER_ERROR", err.Error())
		return
	}
	
	meta := response.MetaInfo{ Total: len(customers), Page: 1, Limit: len(customers) }
	response.SuccessWithMeta(c, 200, customers, meta)
}

func (h *CustomerHandler) Get(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")

	customer, err := h.service.GetCustomerByID(id, tenantID.(string))
	if err != nil {
		response.Error(c, 404, "NOT_FOUND", err.Error())
		return
	}
	response.Success(c, 200, customer)
}
