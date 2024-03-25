package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	myGin "github.com/NeptuneYeh/simplebank/init/gin"
	"github.com/NeptuneYeh/simplebank/init/store"
	mockdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/mock"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/NeptuneYeh/simplebank/tools/hashPassword"
	"github.com/NeptuneYeh/simplebank/tools/helper"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := postgresdb.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		//{
		//	name: "InternalError",
		//	body: gin.H{
		//		"username":  user.Username,
		//		"password":  password,
		//		"full_name": user.FullName,
		//		"email":     user.Email,
		//	},
		//	buildStubs: func(store *mockdb.MockStore) {
		//		store.EXPECT().
		//			CreateUser(gomock.Any(), gomock.Any()).
		//			Times(1).
		//			Return(postgresdb.User{}, sql.ErrConnDone)
		//	},
		//	checkResponse: func(recorder *httptest.ResponseRecorder) {
		//		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		//	},
		//},
		//{
		//	name: "DuplicateUsername",
		//	body: gin.H{
		//		"username":  user.Username,
		//		"password":  password,
		//		"full_name": user.FullName,
		//		"email":     user.Email,
		//	},
		//	buildStubs: func(store *mockdb.MockStore) {
		//		store.EXPECT().
		//			CreateUser(gomock.Any(), gomock.Any()).
		//			Times(1).
		//			Return(postgresdb.User{}, postgresdb.ErrUniqueViolation)
		//	},
		//	checkResponse: func(recorder *httptest.ResponseRecorder) {
		//		require.Equal(t, http.StatusForbidden, recorder.Code)
		//	},
		//},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":  "invalid-user#1",
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     "invalid-email",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"username":  user.Username,
				"password":  "123",
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mockdb.NewMockStore(ctrl)
			tc.buildStubs(mockStore)

			_ = store.NewModuleForTest(mockStore)
			ginModule := myGin.NewModule()
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			ginModule.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomUser(t *testing.T) (user postgresdb.User, password string) {
	randomNumber := rand.Intn(10000)
	password, _ = helper.RandomString(6)
	hashedPassword, err := hashPassword.HashPassword(password)
	require.NoError(t, err)

	user = postgresdb.User{
		Username:       "tom_" + fmt.Sprintf("%04d", randomNumber),
		HashedPassword: hashedPassword,
		FullName:       "tom_" + fmt.Sprintf("%04d", randomNumber),
		Email:          "tom_" + fmt.Sprintf("%04d", randomNumber) + "@yopmail.com",
	}
	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user postgresdb.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser postgresdb.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
