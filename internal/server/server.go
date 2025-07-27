// Package server contains internal server set up and business logic
package server

import (
	"log"
	pbTransaction "money-go/api/v1/transaction"
	pbUser "money-go/api/v1/user"
	"money-go/internal/repository"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	pbTransaction.UnimplementedTransactionServiceServer
	pbUser.UnimplementedUserServiceServer
	transactionRepo *repository.TransactionRepository
	userRepo        *repository.UserRepository
}

func NewGRPCServer(db *sqlx.DB) *grpc.Server {
	log.Println("Starting gRPC server setup...")
	gsrv := grpc.NewServer()
	reflection.Register(gsrv)

	log.Println("Initializing repositories...")
	transactionRepo := repository.NewTransactionRepository(db)
	userRepo := repository.NewUserRepository(db)

	srv := &grpcServer{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
	}

	log.Println("Registering gRPC services...")
	pbTransaction.RegisterTransactionServiceServer(gsrv, srv)
	pbUser.RegisterUserServiceServer(gsrv, srv)

	log.Println("gRPC Server setup completed successfully")
	return gsrv
}
