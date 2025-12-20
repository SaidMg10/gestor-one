package http

import (
	"errors"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/SaidMg10/gestor-one/internal/domain"
	"github.com/SaidMg10/gestor-one/internal/service"
	"github.com/gin-gonic/gin"
)

type IncomeHandler struct {
	svc *service.IncomeService
}

func NewIncomeHandler(svc *service.IncomeService) *IncomeHandler {
	return &IncomeHandler{
		svc: svc,
	}
}

type CreateIncomeRequest struct {
	Amount      float64 `form:"amount" binding:"required"`
	Description string  `form:"description" binding:"required"`
	Type        string  `form:"type" binding:"required"`
}

type UpdateIncomeRequest struct {
	Amount      *float64   `form:"amount"`
	Description *string    `form:"description"`
	Type        *string    `form:"type"`
	Date        *time.Time `form:"date"`
}

type IncomeResponse struct {
	ID          uint    `json:"id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	ReceiptFile string  `json:"receipt_file"`
}

func (h *IncomeHandler) Create(c *gin.Context) {
	var req CreateIncomeRequest
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
			log.Printf("failed to close file: %v", err)
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

	income := &domain.Income{
		Amount:      req.Amount,
		Description: req.Description,
		Type:        domain.IncomeType(req.Type),
		CreatedBy:   user.ID,
	}

	if err := h.svc.Create(c.Request.Context(), income, file, fileHeader); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := IncomeResponse{
		ID:          income.ID,
		Amount:      income.Amount,
		Description: income.Description,
		Type:        string(income.Type),
		ReceiptFile: income.Receipt.FileName,
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *IncomeHandler) List(c *gin.Context) {
	incomes, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Entity mapping
	incomeResponses := make([]IncomeResponse, len(incomes))
	for i, income := range incomes {
		incomeResponses[i] = IncomeResponse{
			ID:          income.ID,
			Amount:      income.Amount,
			Description: income.Description,
			Type:        string(income.Type),
			ReceiptFile: income.Receipt.FileName,
		}
	}
	c.JSON(http.StatusOK, incomeResponses)
}

func (h *IncomeHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid income ID"})
		return
	}
	income, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	response := IncomeResponse{
		ID:          income.ID,
		Amount:      income.Amount,
		Description: income.Description,
		Type:        string(income.Type),
		ReceiptFile: income.Receipt.RelPath,
	}
	c.JSON(http.StatusOK, response)
}

func (h *IncomeHandler) DownloadReceipt(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid income ID"})
		return
	}

	income, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// RelPath: /uploads/archivo.pdf -> archivo real en ./uploads
	filename := filepath.Base(income.Receipt.RelPath)
	fullPath := filepath.Join("./uploads", filename)

	c.FileAttachment(fullPath, filename)
}

func (h *IncomeHandler) Update(c *gin.Context) {
	var req UpdateIncomeRequest

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

	// Convertir a partial domain.Income
	income := &domain.Income{}

	if req.Amount != nil {
		income.Amount = *req.Amount
	}
	if req.Description != nil {
		income.Description = *req.Description
	}
	if req.Type != nil {
		income.Type = domain.IncomeType(*req.Type)
	}
	if req.Date != nil {
		income.Date = *req.Date
	}

	file, fileHeader, _ := c.Request.FormFile("receipt")

	err := h.svc.Update(c, uint(id), income, file, fileHeader, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "income updated"})
}

func (h *IncomeHandler) SoftDelete(c *gin.Context) {
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
			c.JSON(http.StatusNotFound, gin.H{"error": "income not found"})
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "income soft deleted"})
}

func (h *IncomeHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.svc.Delete(c.Request.Context(), uint(id))
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "income not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "income permanently deleted"})
}

func (h *IncomeHandler) Restore(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.svc.Restore(c.Request.Context(), uint(id))
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "income not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "income restored"})
}
