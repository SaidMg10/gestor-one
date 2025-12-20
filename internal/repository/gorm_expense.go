package repository

import (
	"context"
	"time"

	"github.com/SaidMg10/gestor-one/internal/domain"
	"gorm.io/gorm"
)

type GormExpenseRepo struct {
	db *gorm.DB
}

func NewGormExpenseRepo(db *gorm.DB) domain.ExpenseRepo {
	return &GormExpenseRepo{db}
}

func (r *GormExpenseRepo) GetByID(ctx context.Context, id uint) (*domain.Expense, error) {
	var expense domain.Expense
	if err := r.db.WithContext(ctx).
		Preload("Receipt").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&expense, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &expense, nil
}

func (r *GormExpenseRepo) List(ctx context.Context) ([]domain.Expense, error) {
	var expenses []domain.Expense
	if err := r.db.WithContext(ctx).
		Preload("Receipt").
		Where("deleted_at IS NULL").
		Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (r *GormExpenseRepo) CreateWithReceipt(
	ctx context.Context,
	expense *domain.Expense,
	receipt *domain.Receipt,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(expense).Error; err != nil {
			return err
		}
		receipt.ExpenseID = &expense.ID
		receipt.IncomeID = nil
		if err := tx.Create(receipt).Error; err != nil {
			return err
		}
		expense.Receipt = *receipt
		return nil
	})
}

func (r *GormExpenseRepo) UpdateWithReceipt(
	ctx context.Context,
	expense *domain.Expense,
	receipt *domain.Receipt,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&domain.Expense{}).
			Where("id = ? AND deleted_at IS NULL", expense.ID).
			Updates(expense)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return domain.ErrNotFound
		}

		if receipt != nil {
			receipt.ID = expense.Receipt.ID
			receipt.ExpenseID = &expense.ID

			if err := tx.Save(receipt).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *GormExpenseRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&domain.Expense{}, id).Error
}

func (r *GormExpenseRepo) SoftDelete(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&domain.Expense{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now).Error
}

func (r *GormExpenseRepo) Restore(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&domain.Expense{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}
