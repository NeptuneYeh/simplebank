package controllers

import (
	"fmt"
	"github.com/NeptuneYeh/simplebank/init/gin"
	mockdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/mock"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccountAPI(t *testing.T) {
	account := createTestAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).Return(account, nil)

	ginModule := gin.NewModule(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	ginModule.Router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
}

func createTestAccount() postgresdb.Account {
	randomNumber := rand.Intn(10000)
	return postgresdb.Account{
		ID:       1,
		Owner:    "tom_" + fmt.Sprintf("%04d", randomNumber),
		Balance:  100,
		Currency: "USD",
	}
}
