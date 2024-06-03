package gRPC

import (
	"context"
	"fmt"
	"github.com/NeptuneYeh/simplebank/init/auth"
	"github.com/NeptuneYeh/simplebank/init/config"
	"github.com/NeptuneYeh/simplebank/init/logger"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/NeptuneYeh/simplebank/tools/hashPassword"
	"github.com/NeptuneYeh/simplebank/tools/helper"
	"github.com/NeptuneYeh/simplebank/tools/token"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	configModule := config.NewModule("../../")
	configModule.AccessTokenDuration = time.Minute
	_ = logger.NewModule()
	_ = auth.NewModule()
	os.Exit(m.Run())
}

func newContextWithBearerToken(t *testing.T, tokenMaker token.Maker, username string, duration time.Duration) context.Context {
	accessToken, _, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	bearerToken := fmt.Sprintf("%s %s", "bearer", accessToken)
	md := metadata.MD{
		"authorization": []string{
			bearerToken,
		},
	}

	return metadata.NewIncomingContext(context.Background(), md)
}

func randomUser(t *testing.T) (user postgresdb.User, password string) {
	randomNumber := rand.Intn(10000)
	password, _ = helper.RandomString(6)
	hashedPassword, err := hashPassword.HashPassword(password)
	require.NoError(t, err)

	user = postgresdb.User{
		Username:       "tom" + fmt.Sprintf("%04d", randomNumber),
		HashedPassword: hashedPassword,
		FullName:       "tom" + fmt.Sprintf("%04d", randomNumber),
		Email:          "tom" + fmt.Sprintf("%04d", randomNumber) + "@yopmail.com",
	}
	return
}
