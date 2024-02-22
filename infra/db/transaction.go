package db

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ID        string `gorm:"primaryKey"`
	Reference string
	UserID    string
	Amount    int64
	CreatedAt time.Time
}

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (tr *TransactionRepository) AddTransaction(id string, userID, reference string, amount int64) error {
	result := tr.db.Create(&Transaction{ID: id, UserID: userID, Reference: fmt.Sprintf("%s-%s", reference, userID), Amount: amount, CreatedAt: time.Now()})
	return result.Error
}

func (tr *TransactionRepository) AddTransactionWithCtx(ctx context.Context, id string, userID, reference string, amount int64) error {
	result := tr.db.WithContext(ctx).Create(&Transaction{ID: id, UserID: userID, Reference: fmt.Sprintf("%s-%s", reference, userID), Amount: amount, CreatedAt: time.Now()})
	return result.Error
}

func (tr *TransactionRepository) UserTransactions(userID string) []Transaction {
	var ts []Transaction
	tr.db.Where("user_id = ?", userID).Find(&ts)
	return ts
}
