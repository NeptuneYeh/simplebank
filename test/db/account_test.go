package db

import (
	"context"
	"database/sql"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// 在Go語言中，t 是 testing.T 類型的一個參數，用於撰寫測試。這是Go標準庫中 testing 包提供的一個結構體，它包含了用於報告測試失敗或錯誤的方法。
func createTestAccount(t *testing.T) postgresdb.Account {
	user := createTestUser(t)
	arg := postgresdb.CreateAccountParams{
		Owner:    user.Username,
		Balance:  100,
		Currency: "USD",
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createTestAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createTestAccount(t)
	// context.Background() 是 Go 語言中 context 包提供的一個函數，用於創建一個空的、背景（background）的 context.Context 實例
	account1Get, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account1Get)

	require.Equal(t, account1.ID, account1Get.ID)
	require.Equal(t, account1.Owner, account1Get.Owner)
	require.Equal(t, account1.Currency, account1Get.Currency)
	require.Equal(t, account1.Balance, account1Get.Balance)
	// 這個斷言的目的是確保 account1.CreatedAt 和 account1Get.CreatedAt 之間的時間差在一秒內。
	require.WithinDuration(t, account1.CreatedAt, account1Get.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createTestAccount(t)

	arg := postgresdb.UpdateAccountParams{
		ID:      account1.ID,
		Balance: 1000,
	}

	account1Update, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account1Update)

	require.Equal(t, account1.ID, account1Update.ID)
	require.Equal(t, account1.Owner, account1Update.Owner)
	require.Equal(t, account1.Currency, account1Update.Currency)
	require.Equal(t, arg.Balance, account1Update.Balance)
	require.WithinDuration(t, account1.CreatedAt, account1Update.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createTestAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account1Get, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	// 這種斷言用於確保一個錯誤（err）的字串表示與預期的字串相等。
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account1Get)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createTestAccount(t)
	}

	arg := postgresdb.ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
