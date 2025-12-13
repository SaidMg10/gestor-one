package repository

import (
	"context"
	"time"

	"github.com/SaidMg10/gestor-one/internal/domain"
	"gorm.io/gorm"
)

type GormIncomeRepo struct {
	db *gorm.DB
}

func NewGormIncomeRepo(db *gorm.DB) domain.IncomeRepo {
	return &GormIncomeRepo{db}
}

func (r *GormIncomeRepo) GetByID(ctx context.Context, id uint) (*domain.Income, error) {
	var income domain.Income
	if err := r.db.WithContext(ctx).
		Preload("Receipt").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&income, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &income, nil
}

func (r *GormIncomeRepo) List(ctx context.Context) ([]domain.Income, error) {
	var incomes []domain.Income
	if err := r.db.WithContext(ctx).
		Preload("Receipt").
		Where("deleted_at IS NULL").
		Find(&incomes).Error; err != nil {
		return nil, err
	}
	return incomes, nil
}

func (r *GormIncomeRepo) Create(ctx context.Context, income *domain.Income) error {
	return r.db.WithContext(ctx).Create(income).Error
}

func (r *GormIncomeRepo) CreateWithReceipt(ctx context.Context, income *domain.Income, receipt *domain.Receipt) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(income).Error; err != nil {
			return err
		}

		// En el repo, después de tx.Create(income)
		receipt.IncomeID = income.ID
		if err := tx.Create(receipt).Error; err != nil {
			return err
		}
		// Opcionalmente:
		income.Receipt = *receipt

		return nil
	})
}

func (r *GormIncomeRepo) Update(ctx context.Context, income *domain.Income) error {
	return r.db.WithContext(ctx).
		Model(&income).
		Where("id = ?", income.ID).
		Updates(income).
		Error
}

func (r *GormIncomeRepo) UpdateWithReceipt(
	ctx context.Context,
	income *domain.Income,
	receipt *domain.Receipt,
) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&domain.Income{}).
			Where("id = ? AND deleted_at IS NULL", income.ID).
			Updates(income)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return domain.ErrNotFound
		}

		if receipt != nil {
			receipt.ID = income.Receipt.ID
			receipt.IncomeID = income.ID

			if err := tx.Save(receipt).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *GormIncomeRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&domain.Income{}, id).Error
}

func (r *GormIncomeRepo) SoftDelete(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&domain.Income{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now).Error
}

func (r *GormIncomeRepo) Restore(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&domain.Income{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}
