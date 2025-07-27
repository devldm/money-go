package models

import (
	"time"

	pbUser "money-go/api/v1/user"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type User struct {
	ID        string          `db:"id" json:"id"`
	Name      string          `db:"name" json:"name"`
	Email     string          `db:"email" json:"email"`
	Balance   decimal.Decimal `db:"balance" json:"balance"`
	CreatedAt string          `db:"created_at" json:"created_at"`
}

func NewUser(name, email string) *User {
	return &User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Balance:   decimal.Zero,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
}

func (u *User) ToProto() *pbUser.User {
	createdAt, _ := time.Parse(time.RFC3339, u.CreatedAt)
	return &pbUser.User{
		Id:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Balance:   u.Balance.String(),
		CreatedAt: timestamppb.New(createdAt),
	}
}
