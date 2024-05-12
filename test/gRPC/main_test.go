package gRPC

import (
	"context"
	"fmt"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/NeptuneYeh/simplebank/tools/hashPassword"
	"github.com/NeptuneYeh/simplebank/tools/helper"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/test/bufconn"
	"math/rand"
	"net"
	"os"
	"testing"
)

var lis *bufconn.Listener

func TestMain(m *testing.M) {
	os.Exit(m.Run())
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

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}
