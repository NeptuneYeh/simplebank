package hashPassword

import (
	hp "github.com/NeptuneYeh/simplebank/tools/hashPassword"
	"github.com/NeptuneYeh/simplebank/tools/helper"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPassword(t *testing.T) {
	password, err := helper.RandomString(8)
	require.NoError(t, err)
	hashPassword, err := hp.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword)

	err = hp.CheckPassword(password, hashPassword)
	require.NoError(t, err)

	wrongPassword, err := helper.RandomString(6)
	require.NoError(t, err)
	err = hp.CheckPassword(wrongPassword, hashPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashPassword2, err := hp.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword2)
	require.NotEqual(t, hashPassword, hashPassword2)
}
