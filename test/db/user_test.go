package db

import (
	"context"
	"fmt"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func createTestUser(t *testing.T) postgresdb.User {
	randomNumber := rand.Intn(10000)
	arg := postgresdb.CreateUserParams{
		Username:       "tom_" + fmt.Sprintf("%04d", randomNumber),
		HashedPassword: "secret",
		FullName:       "tom_" + fmt.Sprintf("%04d", randomNumber),
		Email:          "tom_" + fmt.Sprintf("%04d", randomNumber) + "@yopmail.com",
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createTestUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createTestUser(t)
	user1Get, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user1Get)

	require.Equal(t, user1.Username, user1Get.Username)
	require.Equal(t, user1.HashedPassword, user1Get.HashedPassword)
	require.Equal(t, user1.FullName, user1Get.FullName)
	require.Equal(t, user1.Email, user1Get.Email)

	require.WithinDuration(t, user1.CreatedAt, user1Get.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user1Get.PasswordChangedAt, time.Second)
}
