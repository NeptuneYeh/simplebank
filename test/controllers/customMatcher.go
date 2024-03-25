package controllers

import (
	"fmt"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/NeptuneYeh/simplebank/tools/hashPassword"
	"github.com/golang/mock/gomock"
	"reflect"
)

type eqCreateUserParamsMatcher struct {
	arg      postgresdb.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(postgresdb.CreateUserParams)
	if !ok {
		return false
	}

	err := hashPassword.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg postgresdb.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}
