package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/SaidMg10/gestor-one/internal/domain"
)

type IncomeService struct {
	incomeRepo  domain.IncomeRepo
	fileStorage domain.FileStorage
}

func NewIncomeService(i domain.IncomeRepo, fS domain.FileStorage) *IncomeService {
	return &IncomeService{
		incomeRepo:  i,
		fileStorage: fS,
	}
}

func (s *IncomeService) Create(
	ctx context.Context,
	income *domain.Income,
	file multipart.File,
	fileHeader *multipart.FileHeader,
) error {
	if income == nil {
		return errors.New("income cannot be nil")
	}
	if income.Amount <= 0 {
		return errors.New("income amount must be greater than 0")
	}
	if income.Description == "" {
		return errors.New("income description is required")
	}
	if income.Type == "" {
		return errors.New("income type is required")
	}
	if !domain.IsValidReceiptType(income.Type) {
		return errors.New("invalid income type")
	}
	if income.CreatedBy == 0 {
		return errors.New("income created_by is required")
	}
	if income.Date.IsZero() {
		income.Date = time.Now()
	}

	fileName, checksum, fullPath, err := s.fileStorage.SavePDF(fileHeader)
	if err != nil {
		return fmt.Errorf("failed to save pdf: %w", err)
	}

	// Crear Receipt con la URL donde se guardó
	receipt := &domain.Receipt{
		FileName:   fileName,
		FileURL:    fullPath,
		MimeType:   "application/pdf",
		UploadedBy: income.CreatedBy,
		Checksum:   checksum,
	}

	err = s.incomeRepo.CreateWithReceipt(ctx, income, receipt)
	if err != nil {
		// Si DB falla, eliminar archivo para evitar basura
		if rmErr := s.fileStorage.DeletePDF(fullPath); rmErr != nil {
			fmt.Printf("failed to remove file after tx error: %v", rmErr)
		}
		return fmt.Errorf("failed to create income with receipt: %w", err)
	}

	return nil
}

func (s *IncomeService) GetByID(ctx context.Context, id uint) (*domain.Income, error) {
	return s.incomeRepo.GetByID(ctx, id)
}

func (s *IncomeService) List(ctx context.Context) ([]domain.Income, error) {
	return s.incomeRepo.List(ctx)
}

func (s *IncomeService) Update(
	ctx context.Context,
	id uint,
	partial *domain.Income,
	file multipart.File,
	fileHeader *multipart.FileHeader,
	userID uint,
) error {
	// Obtener income existente
	existing, err := s.incomeRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrNotFound
	}

	// Verificar permisos
	if existing.CreatedBy != userID {
		return errors.New("only the creator can update this income/receipt")
	}

	// Partial update de campos de Income
	if partial.Amount != 0 {
		existing.Amount = partial.Amount
	}
	if partial.Description != "" {
		existing.Description = partial.Description
	}
	if !partial.Date.IsZero() {
		existing.Date = partial.Date
	}
	if partial.Type != "" {
		existing.Type = partial.Type
	}

	var oldFilePath string
	var receiptToUpdate *domain.Receipt

	// Actualizar receipt solo si hay un archivo nuevo
	if fileHeader != nil {
		if existing.Receipt.ID == 0 {
			return errors.New("receipt not found for this income")
		}

		filename, checksum, fullPath, err := s.fileStorage.SavePDF(fileHeader)
		if err != nil {
			return fmt.Errorf("failed to save PDF: %w", err)
		}

		// Guardar path antiguo para eliminar después si cambia
		if existing.Receipt.Checksum != "" && existing.Receipt.Checksum != checksum {
			oldFilePath = existing.Receipt.FileURL
		}

		// Actualizar campos del receipt existente
		existing.Receipt.FileName = filename
		existing.Receipt.FileURL = fullPath
		existing.Receipt.MimeType = "application/pdf"
		existing.Receipt.UploadedBy = userID
		existing.Receipt.Checksum = checksum

		receiptToUpdate = &existing.Receipt
	}

	// Llamar al repo con receipt actualizado o nil si no hay cambios
	err = s.incomeRepo.UpdateWithReceipt(ctx, existing, receiptToUpdate)
	if err != nil {
		// rollback del archivo nuevo si hubo error
		if receiptToUpdate != nil {
			if rmErr := s.fileStorage.DeletePDF(receiptToUpdate.FileURL); rmErr != nil {
				fmt.Printf("failed to remove new receipt after update error: %v", rmErr)
			}
		}
		return fmt.Errorf("failed to update income with receipt: %w", err)
	}

	// eliminar archivo antiguo si cambió
	if oldFilePath != "" {
		if rmErr := s.fileStorage.DeletePDF(oldFilePath); rmErr != nil {
			fmt.Printf("failed to remove old receipt file: %v", rmErr)
		}
	}

	return nil
}

func (s *IncomeService) SoftDelete(ctx context.Context, id uint, userID uint) error {
	income, err := s.incomeRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if userID != income.CreatedBy {
		fmt.Println("userID:", userID)
		return errors.New("only the creator can delete this income/receipt")
	}
	return s.incomeRepo.SoftDelete(ctx, id)
}

func (s *IncomeService) Delete(ctx context.Context, id uint) error {
	income, err := s.incomeRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if income == nil {
		return domain.ErrNotFound
	}

	if err := s.incomeRepo.Delete(ctx, id); err != nil {
		return err
	}

	if err := s.fileStorage.DeletePDF(income.Receipt.FileURL); err != nil {
		fmt.Printf("failed to remove receipt file after income delete: %v", err)
	}

	return nil
}

func (s *IncomeService) Restore(ctx context.Context, id uint) error {
	income, err := s.incomeRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if income == nil {
		return domain.ErrNotFound
	}

	if err := s.incomeRepo.Restore(ctx, id); err != nil {
		return err
	}

	return nil
}
