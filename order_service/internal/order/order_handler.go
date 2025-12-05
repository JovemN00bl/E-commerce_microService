package order

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	orders := router.Group("/orders")
	{
		orders.GET("/:id", h.handleGetOrderById)
		orders.POST("", h.handleCreateOrder)

	}
}

type createOrderRequest struct {
	UserId string         `json:"user_id" binding:"required"`
	Items  []orderItemDTO `json:"items" binding:"required,dive"`
}

type orderItemDTO struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

func (h *Handler) handleCreateOrder(c *gin.Context) {
	var req createOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	serviceItems := make([]OrderItemInput, len(req.Items))
	for i, items := range req.Items {
		serviceItems[i] = OrderItemInput{
			ProductId: items.ProductID,
			Quantity:  items.Quantity,
		}
	}

	input := CreateOrderInput{UserId: req.UserId, Items: serviceItems}

	order, err := h.service.CreateOrder(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ErrInsufficientStock) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, ErrEmptyCart) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar pedido!"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *Handler) handleGetOrderById(c *gin.Context) {
	id := c.Param("id")

	order, err := h.service.GetById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pedido não encontrado!"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar pedido"})
		return
	}
	c.JSON(http.StatusOK, order)
}
