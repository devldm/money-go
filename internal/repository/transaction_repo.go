// Package repository contains database access layer for users and transactions.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"money-go/internal/models"

	"github.com/jmoiron/sqlx"
)

type TransactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

// CreateTransaction inserts a new transaction into the database
func (r *TransactionRepository) CreateTransaction(ctx context.Context, tx *models.Transaction) error {
	query := `
        INSERT INTO transactions (id, from_user_id, to_user_id, amount, currency, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		tx.ID,
		tx.FromUserID,
		tx.ToUserID,
		tx.Amount,
		tx.Currency,
		tx.Status,
		tx.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

// GetTransactionByID retrieves a transaction by its ID
func (r *TransactionRepository) GetTransactionByID(ctx context.Context, id string) (*models.Transaction, error) {
	var tx models.Transaction

	query := `
        SELECT id, from_user_id, to_user_id, amount, currency, status, created_at
        FROM transactions
        WHERE id = $1`

	err := r.db.GetContext(ctx, &tx, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &tx, nil
}

// GetTransactionsByUserID retrieves all transactions for a user (sent or received)
func (r *TransactionRepository) GetTransactionsByUserID(ctx context.Context, userID string, limit, offset int) ([]*models.Transaction, error) {
	var transactions []*models.Transaction

	query := `
        SELECT id, from_user_id, to_user_id, amount, currency, status, created_at
        FROM transactions
        WHERE from_user_id = $1 OR to_user_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3`

	err := r.db.SelectContext(ctx, &transactions, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
}

// UpdateTransactionStatus updates the status of a transaction
func (r *TransactionRepository) UpdateTransactionStatus(ctx context.Context, id, status string) error {
	query := `
        UPDATE transactions
        SET status = $1
        WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}

// GetTransactionHistory gets paginated transaction history
func (r *TransactionRepository) GetTransactionHistory(ctx context.Context, limit, offset int) ([]*models.Transaction, error) {
	var transactions []*models.Transaction

	query := `
        SELECT id, from_user_id, to_user_id, amount, currency, status, created_at
        FROM transactions
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2`

	err := r.db.SelectContext(ctx, &transactions, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction history: %w", err)
	}

	return transactions, nil
}
