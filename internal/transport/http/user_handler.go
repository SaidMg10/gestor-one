package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/SaidMg10/gestor-one/internal/domain"
	"github.com/SaidMg10/gestor-one/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

// CreateUserRequest represents a request for creating a new user entity in the system.
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
	Active   *bool  `json:"active"`
}

// UpdateUserRequest represents a request for updating a user entity in the system.
type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty"`
	LastName *string `json:"last_name,omitempty"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	Phone    *string `json:"phone,omitempty"`
	Password *string `json:"password,omitempty"`
	Role     *string `json:"role,omitempty"`
	Active   *bool   `json:"active"`
}

// UserResponse represents a response for a user entity in the system.
type UserResponse struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	LastName string  `json:"last_name"`
	Email    string  `json:"email"`
	Phone    string  `json:"phone"`
	Role     string  `json:"role"`
	GoogleID *string `json:"google_id,omitempty"`
	Active   bool    `json:"active"`
}

func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userEntity := &domain.User{
		Name:     req.Name,
		LastName: req.LastName,
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     req.Role,
	}

	if err := h.svc.Create(c.Request.Context(), userEntity, req.Password); err != nil {
		// Default status
		status := http.StatusBadRequest

		// Mapear errores específicos a HTTP
		switch {
		case errors.Is(err, domain.ErrEmailExists):
			status = http.StatusConflict
		case strings.Contains(err.Error(), "contraseña"):
			status = http.StatusBadRequest
		case strings.Contains(err.Error(), "email inválido"):
			status = http.StatusBadRequest
		case strings.Contains(err.Error(), "superadmin"):
			status = http.StatusForbidden
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, userEntity)
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Entity mapping
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			LastName: user.LastName,
			Email:    user.Email,
			Phone:    user.Phone,
			Role:     user.Role,
			Active:   *user.Active,
		}
	}
	c.JSON(http.StatusOK, userResponses)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	user, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	response := UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		LastName: user.LastName,
		Email:    user.Email,
		Phone:    user.Phone,
		Role:     user.Role,
		Active:   *user.Active,
	}
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.Active != nil {
		user.Active = req.Active
	}

	var pwd *string
	if req.Password != nil {
		pwd = req.Password
	}

	err = h.svc.Update(c.Request.Context(), uint(id), user, pwd)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmailExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		case errors.Is(err, domain.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	err = h.svc.Delete(c.Request.Context(), uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, domain.ErrNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
