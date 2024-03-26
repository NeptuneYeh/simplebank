package controllers

import (
	"database/sql"
	"github.com/NeptuneYeh/simplebank/init/auth"
	"github.com/NeptuneYeh/simplebank/init/config"
	"github.com/NeptuneYeh/simplebank/init/logger"
	"github.com/NeptuneYeh/simplebank/init/store"
	"github.com/NeptuneYeh/simplebank/internal/application/base"
	"github.com/NeptuneYeh/simplebank/internal/application/requests/userRequests"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/NeptuneYeh/simplebank/tools/hashPassword"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
)

type UserController struct {
	store postgresdb.Store
}

func NewUserController() *UserController {
	return &UserController{
		store: store.MainStore,
	}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var req userRequests.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, base.ErrorResponse(err))
		return
	}

	hashedPassword, err := hashPassword.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, base.ErrorResponse(err))
		return
	}

	arg := postgresdb.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := c.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, base.ErrorResponse(err))
				return
			}
		}
		logger.MainLog.Error().Msg(err.Error())
		ctx.JSON(http.StatusInternalServerError, base.ErrorResponse(err))
		return
	}

	resp := userRequests.UserResponse{
		Username:          user.Username,
		Email:             user.Email,
		FullName:          user.FullName,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *UserController) LoginUser(ctx *gin.Context) {
	var req userRequests.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, base.ErrorResponse(err))
		return
	}

	user, err := c.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, base.ErrorResponse(userRequests.ErrEmailOrPasswordNotCorrect))
			return
		}
		ctx.JSON(http.StatusInternalServerError, base.ErrorResponse(err))
		return
	}

	// check password
	err = hashPassword.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, base.ErrorResponse(userRequests.ErrEmailOrPasswordNotCorrect))
		return
	}

	// create accessToken
	accessToken, err := auth.MainAuth.CreateToken(user.Username, config.MainConfig.AccessTokenDuration)

	resp := userRequests.LoginUserResponse{
		AccessToken: accessToken,
		User:        userRequests.NewUserResponse(user),
	}

	ctx.JSON(http.StatusOK, resp)
}
