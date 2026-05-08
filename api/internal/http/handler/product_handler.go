package handler

import (
	"github.com/gin-gonic/gin"
	"lakoo/backend/internal/dto"
	"lakoo/backend/internal/service"
	"lakoo/backend/pkg/response"
	"log"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(uu service.ProductService) *ProductHandler {
	return &ProductHandler{service: uu}
}

func (h *ProductHandler) Create(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.Error(c, 403, "FORBIDDEN", "Tenant Context Missing")
		return
	}

	var req dto.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "BAD_REQUEST", "Invalid input format")
		return
	}

	userID, _ := c.Get("user_id")

	res, err := h.service.CreateProduct(tenantID.(string), userID.(string), &req)
	if err != nil {
		log.Printf("[ProductHandler] Create Error for Tenant %s: %v", tenantID.(string), err)
		response.Error(c, 422, "DATABASE_ERROR", "Gagal menyimpan produk: "+err.Error())
		return
	}

	response.Success(c, 201, res)
}

func (h *ProductHandler) Update(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")

	var req dto.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "BAD_REQUEST", "Invalid input format")
		return
	}

	userID, _ := c.Get("user_id")

	res, err := h.service.UpdateProduct(id, tenantID.(string), userID.(string), &req)
	if err != nil {
		if err.Error() == "product not found" {
			response.Error(c, 404, "NOT_FOUND", "Product not found")
		} else {
			response.Error(c, 400, "BAD_REQUEST", err.Error())
		}
		return
	}
	response.Success(c, 200, res)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")

	err := h.service.DeleteProduct(id, tenantID.(string))
	if err != nil {
		response.Error(c, 404, "NOT_FOUND", "Product not found or unauthorized")
		return
	}
	response.Success(c, 200, gin.H{"message": "Product successfully deleted"})
}

func (h *ProductHandler) List(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")

	products, err := h.service.GetProducts(tenantID.(string))
	if err != nil {
		response.Error(c, 500, "INTERNAL_SERVER_ERROR", err.Error())
		return
	}
	
	meta := response.MetaInfo{ Total: len(products), Page: 1, Limit: len(products) }
	response.SuccessWithMeta(c, 200, products, meta)
}

func (h *ProductHandler) Get(c *gin.Context) {
	tenantID, _ := c.Get("tenant_id")
	id := c.Param("id")

	product, err := h.service.GetProductByID(id, tenantID.(string))
	if err != nil {
		response.Error(c, 404, "NOT_FOUND", err.Error())
		return
	}
	response.Success(c, 200, product)
}
