package gRPC

import (
	"context"
	"github.com/NeptuneYeh/simplebank/init/store"
	"github.com/NeptuneYeh/simplebank/internal/grpcApp"
	mockdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/mock"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/NeptuneYeh/simplebank/pb"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"testing"
)

const bufSize = 1024 * 1024

func TestGRPCCreateUserAPI(t *testing.T) {
	// TODO mock store
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStore := mockdb.NewMockStore(ctrl)

	// TODO set mockStore to storeModule
	_ = store.NewModuleForTest(mockStore)

	// TODO gRPC server with mockStore
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterSimpleBankServer(s, grpcApp.NewModule())

	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// TODO gRPC client
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewSimpleBankClient(conn)

	// TODO Test detail

	user, password := randomUser(t)

	testCases := []struct {
		name          string
		req           *pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore, req *pb.CreateUserRequest)
		checkResponse func(response *pb.CreateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, req *pb.CreateUserRequest) {
				arg := postgresdb.CreateUserParams{
					Username: req.GetUsername(),
					FullName: req.GetFullName(),
					Email:    req.GetEmail(),
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(response *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, response)
			},
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockStore, tc.req)
			resp, err := client.CreateUser(ctx, tc.req)
			tc.checkResponse(resp, err)
		})
	}
}
