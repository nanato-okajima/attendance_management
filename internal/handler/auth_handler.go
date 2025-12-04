package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/nanato-okajima/attendance_management/internal/entity"
	"github.com/nanato-okajima/attendance_management/internal/errors"
	"github.com/nanato-okajima/attendance_management/internal/service"
)

type AuthHandler interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

type authHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	EmployeeNumber int    `json:"employee_number" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=8"`
	Role           string `json:"role" binding:"required,oneof=admin employee"`
}

func (h *authHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"token": token,
		"type":  "Bearer",
		"user": gin.H{
			"id":              user.ID,
			"employee_number": user.EmployeeNumber,
			"email":           user.Email,
			"role":            user.Role,
		},
	})
}

func (h *authHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user := &entity.User{
		EmployeeNumber: req.EmployeeNumber,
		Email:          req.Email,
		Role:           req.Role,
	}

	if err := h.authService.Register(user, req.Password); err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(201, gin.H{"message": "user registered successfully"})
}
