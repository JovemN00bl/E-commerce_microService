package product

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	productRoutes := router.Group("/products")
	{
		productRoutes.POST("", h.handleCreateProduct)
		productRoutes.GET("/:id", h.handleGetByID)
		productRoutes.GET("", h.handleListProducts)
	}
}

type createProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required, gt=0"`
	Stock       int     `json:"stock" binding:"required, gt=0"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) handleCreateProduct(c *gin.Context) {
	var req createProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	product, err := h.service.Create(
		c.Request.Context(),
		req.Name,
		req.Description,
		req.Price,
		req.Stock)

	if err != nil {
		if errors.Is(err, ErrInvalidPrice) || errors.Is(err, ErrProductNameRequired) || errors.Is(err, ErrInvalidStock) {
			c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Erro ao criar produto."})
		return
	}
	c.JSON(http.StatusCreated, product)
}

func (h *Handler) handleGetByID(c *gin.Context) {
	productID, _ := strconv.Atoi(c.Param("id"))

	product, err := h.service.GetById(c.Request.Context(), productID)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			c.JSON(http.StatusNotFound, errorResponse{Error: "Produto não encontrado."})
			return
		}

		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Erro ao buscar produto"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *Handler) handleListProducts(c *gin.Context) {
	product, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Erro ao listar produtos."})
		return
	}

	c.JSON(http.StatusOK, product)
}
