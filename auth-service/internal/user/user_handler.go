package user

import (
	"errors"
	"net/http"
	_ "syscall"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/register", h.handleRegister)
		userRoutes.POST("/login", h.handleLogin)
	}
}

type registerRequest struct {
	Email    string `json:"Email" binding:"required, email"`
	Password string `json:"password" binding:"required"`
}

type loginRequest struct {
	Email    string `json:"Email" binding:"required, email"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) handleRegister(c *gin.Context) {
	var req registerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	err := h.service.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Error ao criar o usuario"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuario criado com sucesso"})
}

func (h *Handler) handleLogin(c *gin.Context) {
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	token, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, errorResponse{Error: "email ou senha invalidos"})
			return
		}

		c.JSON(http.StatusInternalServerError, errorResponse{Error: "error interno no servidor"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{Token: token})

}
