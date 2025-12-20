package http

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/SaidMg10/gestor-one/internal/domain"
	"github.com/SaidMg10/gestor-one/internal/service"
	"github.com/gin-gonic/gin"
)

type ExpenseHandler struct {
	svc *service.ExpenseService
}

func NewExpenseHandler(svc *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{
		svc: svc,
	}
}

type CreateExpenseRequest struct {
	Amount      float64 `form:"amount" binding:"required"`
	Description string  `form:"description" binding:"required"`
	Type        string  `form:"type" binding:"required"`
}

type UpdateExpenseRequest struct {
	Amount      *float64   `form:"amount"`
	Description *string    `form:"description"`
	Type        *string    `form:"type"`
	Date        *time.Time `form:"date"`
}

type ExpenseResponse struct {
	ID          uint    `json:"id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	ReceiptFile string  `json:"receipt_file"`
}

func (h *ExpenseHandler) Create(c *gin.Context) {
	var req CreateExpenseRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, fileHeader, err := c.Request.FormFile("receipt")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "receipt file is required"})
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("failed to close file: %v\n", err)
		}
	}()

	if filepath.Ext(fileHeader.Filename) != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only PDF files are allowed"})
		return
	}

	userCtx, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, ok := userCtx.(*domain.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user type in context"})
		return
	}

	expense := &domain.Expense{
		Amount:      req.Amount,
		Description: req.Description,
		Type:        domain.ExpenseType(req.Type),
		CreatedBy:   user.ID,
	}

	if err := h.svc.Create(c.Request.Context(), expense, file, fileHeader); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := ExpenseResponse{
		ID:          expense.ID,
		Amount:      expense.Amount,
		Description: expense.Description,
		Type:        string(expense.Type),
		ReceiptFile: expense.Receipt.FileURL,
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *ExpenseHandler) List(c *gin.Context) {
	expenses, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	expenseResponses := make([]ExpenseResponse, len(expenses))
	for i, expense := range expenses {
		expenseResponses[i] = ExpenseResponse{
			ID:          expense.ID,
			Amount:      expense.Amount,
			Description: expense.Description,
			Type:        string(expense.Type),
			ReceiptFile: expense.Receipt.FileURL,
		}
	}
	c.JSON(http.StatusOK, expenseResponses)
}

func (h *ExpenseHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expense ID"})
		return
	}
	expense, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	response := ExpenseResponse{
		ID:          expense.ID,
		Amount:      expense.Amount,
		Description: expense.Description,
		Type:        string(expense.Type),
		ReceiptFile: expense.Receipt.FileURL,
	}
	c.JSON(http.StatusOK, response)
}

func (h *ExpenseHandler) Update(c *gin.Context) {
	var req UpdateExpenseRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	userCtx, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, ok := userCtx.(*domain.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user type in context"})
		return
	}

	expense := &domain.Expense{}

	if req.Amount != nil {
		expense.Amount = *req.Amount
	}
	if req.Description != nil {
		expense.Description = *req.Description
	}
	if req.Type != nil {
		expense.Type = domain.ExpenseType(*req.Type)
	}
	if req.Date != nil {
		expense.Date = *req.Date
	}

	file, fileHeader, _ := c.Request.FormFile("receipt")

	err := h.svc.Update(c, uint(id), expense, file, fileHeader, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "expense updated"})
}

func (h *ExpenseHandler) SoftDelete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userCtx, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, ok := userCtx.(*domain.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user type in context"})
		return
	}
	err = h.svc.SoftDelete(c.Request.Context(), uint(id), user.ID)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "expense not found"})
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "expense soft deleted"})
}

func (h *ExpenseHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.svc.Delete(c.Request.Context(), uint(id))
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "expense not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "expense permanently deleted"})
}

func (h *ExpenseHandler) Restore(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.svc.Restore(c.Request.Context(), uint(id))
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "expense not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "expense restored"})
}
