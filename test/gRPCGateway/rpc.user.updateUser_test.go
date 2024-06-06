package gRPC

import (
	"context"
	"github.com/NeptuneYeh/simplebank/init/asynq"
	"github.com/NeptuneYeh/simplebank/init/auth"
	"github.com/NeptuneYeh/simplebank/init/grpcGateway"
	"github.com/NeptuneYeh/simplebank/init/store"
	mockdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/mock"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/NeptuneYeh/simplebank/pb"
	"github.com/NeptuneYeh/simplebank/tools/token"
	mockwk "github.com/NeptuneYeh/simplebank/tools/worker/mock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGRPCGatewayUpdateUserAPI(t *testing.T) {
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
	user, _ := randomUser(t)

	newName := "HelloAlex"
	newEmail := "helloalex@yopmail.com"

	testCases := []struct {
		name          string
		req           *pb.UpdateUserRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, response *pb.UpdateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := postgresdb.UpdateUserParams{
					Username: user.Username,
					FullName: pgtype.Text{
						String: newName,
						Valid:  true,
					},
					Email: pgtype.Text{
						String: newEmail,
						Valid:  true,
					},
				}
				updatedUser := postgresdb.User{
					Username:          user.Username,
					HashedPassword:    user.HashedPassword,
					FullName:          newName,
					Email:             newEmail,
					PasswordChangedAt: user.PasswordChangedAt,
					CreatedAt:         user.CreatedAt,
					IsEmailVerified:   user.IsEmailVerified,
				}
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updatedUser, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, response *pb.UpdateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, response)
				updatedUser := response.GetUser()
				require.Equal(t, user.Username, updatedUser.Username)
				require.Equal(t, newName, updatedUser.FullName)
				require.Equal(t, newEmail, updatedUser.Email)
			},
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mockStore)
			// TODO set up test server
			grpcGatewayModule := grpcGateway.NewModule()
			ctx := tc.buildContext(t, auth.MainAuth)
			resp, err := grpcGatewayModule.GrpcApi.UpdateUser(ctx, tc.req)
			tc.checkResponse(t, resp, err)
		})
	}
}
