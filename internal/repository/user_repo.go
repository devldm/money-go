package repository

import (
	"context"
	"database/sql"
	"fmt"
	"money-go/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var tx models.User

	query := `
        SELECT * 
        FROM users 
        WHERE id = $1`

	err := r.db.GetContext(ctx, &tx, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &tx, nil
}

// UpdateBalance updates a users balance
func (r *UserRepository) UpdateBalance(ctx context.Context, userID string, newBalance decimal.Decimal) error {
	query := `
        UPDATE users 
        SET balance = $1 
        WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, newBalance.String(), userID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
