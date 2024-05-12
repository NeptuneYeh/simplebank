package grpcApp

import (
	"context"
	"github.com/NeptuneYeh/simplebank/init/logger"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/NeptuneYeh/simplebank/pb"
	"github.com/NeptuneYeh/simplebank/tools/hashPassword"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Module) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := hashPassword.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	arg := postgresdb.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := c.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %v", err)
			}
		}
		logger.MainLog.Error().Msg(err.Error())
		return nil, status.Errorf(codes.AlreadyExists, "failed to create user: %v", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}
