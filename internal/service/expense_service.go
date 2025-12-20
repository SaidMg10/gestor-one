package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/SaidMg10/gestor-one/internal/domain"
)

type ExpenseService struct {
	expenseRepo domain.ExpenseRepo
	fileStorage domain.FileStorage
}

func NewExpenseService(e domain.ExpenseRepo, fS domain.FileStorage) *ExpenseService {
	return &ExpenseService{
		expenseRepo: e,
		fileStorage: fS,
	}
}

func (s *ExpenseService) Create(
	ctx context.Context,
	expense *domain.Expense,
	file multipart.File,
	fileHeader *multipart.FileHeader,
) error {
	if expense == nil {
		return errors.New("expense cannot be nil")
	}
	if expense.Amount <= 0 {
		return errors.New("expense amount must be greater than 0")
	}
	if expense.Description == "" {
		return errors.New("expense description is required")
	}
	if expense.Type == "" {
		return errors.New("expense type is required")
	}
	if !domain.IsValidExpenseType(expense.Type) {
		return errors.New("invalid expense type")
	}
	if expense.CreatedBy == 0 {
		return errors.New("expense created_by is required")
	}
	if expense.Date.IsZero() {
		expense.Date = time.Now()
	}

	fileName, checksum, relPath, err := s.fileStorage.SavePDF(fileHeader)
	if err != nil {
		return fmt.Errorf("failed to save pdf: %w", err)
	}

	receipt := &domain.Receipt{
		FileName:   fileName,
		RelPath:    relPath,
		MimeType:   "application/pdf",
		UploadedBy: expense.CreatedBy,
		Checksum:   checksum,
	}

	err = s.expenseRepo.CreateWithReceipt(ctx, expense, receipt)
	if err != nil {
		if rmErr := s.fileStorage.DeletePDF(relPath); rmErr != nil {
			fmt.Printf("failded to remove file after tx error: %v", rmErr)
		}
		return fmt.Errorf("failed to create expense with receipt: %w", err)
	}
	return nil
}

func (s *ExpenseService) GetByID(ctx context.Context, id uint) (*domain.Expense, error) {
	return s.expenseRepo.GetByID(ctx, id)
}

func (s *ExpenseService) List(ctx context.Context) ([]domain.Expense, error) {
	return s.expenseRepo.List(ctx)
}

func (s ExpenseService) Update(
	ctx context.Context,
	id uint,
	partial *domain.Expense,
	file multipart.File,
	fileHeader *multipart.FileHeader,
	userID uint,
) error {
	existing, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrNotFound
	}
	if existing.CreatedBy != userID {
		return errors.New("only the creator can update this expense/receipt")
	}

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

	if fileHeader != nil {
		if existing.Receipt.ID == 0 {
			return errors.New("receipt not found for this expense")
		}

		fileName, checksum, relPath, err := s.fileStorage.SavePDF(fileHeader)
		if err != nil {
			return fmt.Errorf("failed to save pdf: %w", err)
		}

		if existing.Receipt.Checksum != "" && existing.Receipt.Checksum != checksum {
			oldFilePath = existing.Receipt.RelPath
		}

		existing.Receipt.FileName = fileName
		existing.Receipt.RelPath = relPath
		existing.Receipt.MimeType = "application/pdf"
		existing.Receipt.UploadedBy = userID
		existing.Receipt.Checksum = checksum

		receiptToUpdate = &existing.Receipt
	}

	err = s.expenseRepo.UpdateWithReceipt(ctx, existing, receiptToUpdate)
	if err != nil {
		if receiptToUpdate != nil {
			if rmErr := s.fileStorage.DeletePDF(receiptToUpdate.RelPath); rmErr != nil {
				fmt.Printf("failed to remove new receipt after update error: %v", rmErr)
			}
		}
		return fmt.Errorf("failed to update expense with receipt: %w", err)
	}

	if oldFilePath != "" {
		if rmErr := s.fileStorage.DeletePDF(oldFilePath); rmErr != nil {
			fmt.Printf("failed to remove old receipt file: %v", rmErr)
		}
	}
	return nil
}

func (s *ExpenseService) SoftDelete(ctx context.Context, id uint, userID uint) error {
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if userID != expense.CreatedBy {
		fmt.Println("userID:", userID)
		return errors.New("only the creator can delete this expense/receipt")
	}
	return s.expenseRepo.SoftDelete(ctx, id)
}

func (s *ExpenseService) Delete(ctx context.Context, id uint) error {
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if expense == nil {
		return domain.ErrNotFound
	}

	if err := s.expenseRepo.Delete(ctx, id); err != nil {
		return err
	}

	if err := s.fileStorage.DeletePDF(expense.Receipt.RelPath); err != nil {
		fmt.Printf("failed to remove receipt file after expense delete: %v", err)
	}

	return nil
}

func (s *ExpenseService) Restore(ctx context.Context, id uint) error {
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if expense == nil {
		return domain.ErrNotFound
	}

	if err := s.expenseRepo.Restore(ctx, id); err != nil {
		return err
	}

	return nil
}
