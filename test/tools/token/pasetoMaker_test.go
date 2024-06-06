package token

import (
	"fmt"
	"github.com/NeptuneYeh/simplebank/tools/helper"
	myRole "github.com/NeptuneYeh/simplebank/tools/role"
	myToken "github.com/NeptuneYeh/simplebank/tools/token"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {
	randomString, err := helper.RandomString(32)
	require.NoError(t, err)
	maker, err := myToken.NewPasetoMaker(randomString)
	require.NoError(t, err)

	randomNumber := rand.Intn(10000)
	username := "test_" + fmt.Sprintf("%04d", randomNumber)
	role := myRole.Depositor
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.Equal(t, role, payload.Role)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	randomString, err := helper.RandomString(32)
	require.NoError(t, err)
	maker, err := myToken.NewPasetoMaker(randomString)
	require.NoError(t, err)

	randomNumber := rand.Intn(10000)
	role := myRole.Depositor
	username := "test_" + fmt.Sprintf("%04d", randomNumber)

	duration := -time.Minute
	token, payload, err := maker.CreateToken(username, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, myToken.ErrExpiredToken.Error())
	require.Nil(t, payload)
}
