package server

import (
	"context"
	"log"
	pbTransaction "money-go/api/v1/transaction"
	"money-go/internal/models"
	"time"

	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *grpcServer) SendMoney(ctx context.Context, in *pbTransaction.SendMoneyRequest) (*pbTransaction.Transaction, error) {
	if in.FromUserId == "" || in.ToUserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user IDs cannot be empty")
	}

	if in.FromUserId == in.ToUserId {
		return nil, status.Error(codes.InvalidArgument, "cannot send money to yourself")
	}

	sender, err := s.userRepo.GetUserByID(ctx, in.FromUserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "sender not found")
	}

	receiver, err := s.userRepo.GetUserByID(ctx, in.ToUserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "receiver not found")
	}

	amount, err := decimal.NewFromString(in.Amount)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid amount")
	}

	if sender.Balance.LessThan(amount) {
		return nil, status.Error(codes.FailedPrecondition, "insufficient balance")
	}

	transaction := models.NewTransaction(in.FromUserId, in.ToUserId, amount, in.Currency)

	err = s.transactionRepo.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create transaction")
	}

	newSenderBalance := sender.Balance.Sub(amount)
	newReceiverBalance := receiver.Balance.Add(amount)

	err = s.userRepo.UpdateBalance(ctx, in.FromUserId, newSenderBalance)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update sender balance")
	}

	err = s.userRepo.UpdateBalance(ctx, in.ToUserId, newReceiverBalance)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update receiver balance")
	}

	createdAtTime, err := time.Parse(time.RFC3339, transaction.CreatedAt)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid timestamp format")
	}

	return &pbTransaction.Transaction{
		Id:         transaction.ID,
		FromUserId: transaction.FromUserID,
		ToUserId:   transaction.ToUserID,
		Amount:     transaction.Amount.String(),
		Currency:   transaction.Currency,
		Status:     transaction.Status,
		CreatedAt:  timestamppb.New(createdAtTime),
	}, nil
}

func (s *grpcServer) GetTransaction(ctx context.Context, in *pbTransaction.GetTransactionRequest) (*pbTransaction.Transaction, error) {
	log.Printf("GetTransaction called: %s", in.TransactionId)

	transaction, err := s.transactionRepo.GetTransactionByID(ctx, in.TransactionId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "transaction not found")
	}

	createdAtTime, err := time.Parse(time.RFC3339, transaction.CreatedAt)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid timestamp format")
	}

	return &pbTransaction.Transaction{
		Id:         in.TransactionId,
		FromUserId: transaction.FromUserID,
		ToUserId:   transaction.ToUserID,
		Amount:     transaction.Amount.String(),
		Currency:   transaction.Currency,
		Status:     transaction.Status,
		CreatedAt:  timestamppb.New(createdAtTime),
	}, nil
}
