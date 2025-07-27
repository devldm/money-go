// Package models contains domain entities for users and transactions.
package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID         string          `db:"id" json:"id"`
	FromUserID string          `db:"from_user_id" json:"from_user_id"`
	ToUserID   string          `db:"to_user_id" json:"to_user_id"`
	Amount     decimal.Decimal `db:"amount" json:"amount"`
	Currency   string          `db:"currency" json:"currency"`
	Status     string          `db:"status" json:"status"`
	CreatedAt  string          `db:"created_at" json:"created_at"`
}

func NewTransaction(fromUserID, toUserID string, amount decimal.Decimal, currency string) *Transaction {
	return &Transaction{
		ID:         uuid.New().String(),
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Amount:     amount,
		Currency:   currency,
		Status:     "pending",
		CreatedAt:  time.Now().Format(time.RFC3339),
	}
}
