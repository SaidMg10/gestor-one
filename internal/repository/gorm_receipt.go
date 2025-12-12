package repository

import (
	"context"

	"github.com/SaidMg10/gestor-one/internal/domain"
	"gorm.io/gorm"
)

type GormReceiptRepo struct {
	db *gorm.DB
}

func NewGormReceiptRepo(db *gorm.DB) domain.ReceiptRepo {
	return &GormReceiptRepo{db: db}
}

func (r *GormReceiptRepo) WithTx(tx *gorm.DB) domain.ReceiptRepo {
	return &GormReceiptRepo{db: tx}
}

func (r *GormReceiptRepo) GetByID(ctx context.Context, id uint) (*domain.Receipt, error) {
	var receipt domain.Receipt
	if err := r.db.WithContext(ctx).First(&receipt, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err

	}
	return &receipt, nil
}

func (r *GormReceiptRepo) List(ctx context.Context) ([]domain.Receipt, error) {
	var receipts []domain.Receipt
	if err := r.db.WithContext(ctx).Find(&receipts).Error; err != nil {
		return nil, err
	}
	return receipts, nil
}

func (r *GormReceiptRepo) Create(ctx context.Context, receipt *domain.Receipt) error {
	return r.db.WithContext(ctx).Create(receipt).Error
}

func (r *GormReceiptRepo) Update(ctx context.Context, receipt *domain.Receipt) error {
	return r.db.WithContext(ctx).
		Model(&domain.Receipt{}).
		Where("id = ?", receipt.ID).
		Updates(receipt).
		Error
}

func (r *GormReceiptRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Receipt{}, id).Error
}
