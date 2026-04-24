package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/greg901896/go-shopflow/internal/model"
	"github.com/greg901896/go-shopflow/internal/service"
)

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

type createProductRequest struct {
	Name  string `json:"name" binding:"required"`
	Price string `json:"price" binding:"required"`
	Stock int    `json:"stock" binding:"gte=0"`
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p := &model.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}
	if err := h.svc.Create(c.Request.Context(), p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}
