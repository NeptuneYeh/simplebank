package gRPC

import (
	"context"
	"github.com/NeptuneYeh/simplebank/init/asynq"
	"github.com/NeptuneYeh/simplebank/init/grpcGateway"
	"github.com/NeptuneYeh/simplebank/init/store"
	mockdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/mock"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/NeptuneYeh/simplebank/pb"
	"github.com/NeptuneYeh/simplebank/tools/worker"
	mockwk "github.com/NeptuneYeh/simplebank/tools/worker/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGRPCGatewayCreateUserAPI(t *testing.T) {
	// TODO mock store
	ctrl1 := gomock.NewController(t)
	defer ctrl1.Finish()
	mockStore := mockdb.NewMockStore(ctrl1)
	// TODO set mockStore to storeModule
	_ = store.NewModuleForTest(mockStore)

	// TODO mock worker
	ctrl2 := gomock.NewController(t)
	defer ctrl2.Finish()
	mockTaskDistributor := mockwk.NewMockTaskDistributor(ctrl2)
	// TODO set mockTaskDistributor to storeModule
	_ = asynq.NewModuleForTest(mockTaskDistributor)

	// TODO Test detail
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		req           *pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore, mockTaskDistributor *mockwk.MockTaskDistributor)
		checkResponse func(t *testing.T, response *pb.CreateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, mockTaskDistributor *mockwk.MockTaskDistributor) {
				arg := postgresdb.CreateUserTxParams{
					CreateUserParams: postgresdb.CreateUserParams{
						Username: user.Username,
						FullName: user.FullName,
						Email:    user.Email,
					},
				}
				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, password, user)).
					Times(1).
					Return(postgresdb.CreateUserTxResult{User: user}, nil)

				taskPayload := &worker.PayloadSendVerifyEmail{
					Username: user.Username,
				}
				mockTaskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, response *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, response)
				createdUser := response.GetUser()
				require.Equal(t, user.Username, createdUser.Username)
				require.Equal(t, user.FullName, createdUser.FullName)
				require.Equal(t, user.Email, createdUser.Email)
			},
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockStore, mockTaskDistributor)
			// TODO set up test server
			grpcGatewayModule := grpcGateway.NewModule()
			resp, err := grpcGatewayModule.GrpcApi.CreateUser(context.Background(), tc.req)
			tc.checkResponse(t, resp, err)
		})
	}
}
