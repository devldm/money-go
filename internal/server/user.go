package server

import (
	"context"
	pbUser "money-go/api/v1/user"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *grpcServer) GetUser(ctx context.Context, in *pbUser.GetUserRequest) (*pbUser.User, error) {
	if in.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID cannot be empty")
	}

	user, err := s.userRepo.GetUserByID(ctx, in.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	createdAtTime, err := time.Parse(time.RFC3339, user.ToProto().CreatedAt.AsTime().Format(time.RFC3339))
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid timestamp format")
	}

	return &pbUser.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Balance:   user.Balance.String(),
		CreatedAt: timestamppb.New(createdAtTime),
	}, nil
}
