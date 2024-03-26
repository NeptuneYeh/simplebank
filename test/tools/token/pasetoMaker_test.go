package token

import (
	"fmt"
	"github.com/NeptuneYeh/simplebank/tools/helper"
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

	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	randomString, err := helper.RandomString(32)
	require.NoError(t, err)
	maker, err := myToken.NewPasetoMaker(randomString)
	require.NoError(t, err)

	randomNumber := rand.Intn(10000)
	username := "test_" + fmt.Sprintf("%04d", randomNumber)

	duration := -time.Minute
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, myToken.ErrExpiredToken.Error())
	require.Nil(t, payload)
}
