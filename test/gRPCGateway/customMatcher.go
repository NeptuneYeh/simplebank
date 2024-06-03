package gRPC

import (
	"fmt"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/NeptuneYeh/simplebank/tools/hashPassword"
	"github.com/golang/mock/gomock"
	"reflect"
)

type eqCreateUserTxParamsMatcher struct {
	arg      postgresdb.CreateUserTxParams
	password string
	user     postgresdb.User
}

func (expected eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	//fmt.Println(">> check param matches")
	actualArg, ok := x.(postgresdb.CreateUserTxParams)
	if !ok {
		return false
	}

	//fmt.Println(">> check password, actualArg.")
	err := hashPassword.CheckPassword(expected.password, actualArg.HashedPassword)
	if err != nil {
		return false
	}

	//fmt.Println(">> deep equal.")
	expected.arg.HashedPassword = actualArg.HashedPassword
	if !reflect.DeepEqual(expected.arg.CreateUserParams, actualArg.CreateUserParams) {
		return false
	}

	// TODO 這裡是一個適合執行 AfterCreate 的地方
	err = actualArg.AfterCreate(expected.user)

	return err == nil
}

func (expected eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", expected.arg, expected.password)
}

func EqCreateUserTxParams(arg postgresdb.CreateUserTxParams, password string, user postgresdb.User) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg, password, user}
}
