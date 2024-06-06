package controllers

import (
	"errors"
	"fmt"
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
	"time"
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
		if errors.Is(err, postgresdb.ErrRecordNotFound) {
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
	accessToken, accessPayload, err := auth.MainAuth.CreateToken(user.Username, config.MainConfig.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, base.ErrorResponse(userRequests.ErrEmailOrPasswordNotCorrect))
		return
	}

	// create refreshToken
	refreshToken, refreshPayload, err := auth.MainAuth.CreateToken(user.Username, config.MainConfig.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, base.ErrorResponse(userRequests.ErrEmailOrPasswordNotCorrect))
		return
	}

	session, err := c.store.CreateSession(ctx, postgresdb.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, base.ErrorResponse(userRequests.ErrEmailOrPasswordNotCorrect))
		return
	}

	resp := userRequests.LoginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  userRequests.NewUserResponse(user),
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *UserController) RenewAccessToken(ctx *gin.Context) {
	var req userRequests.RenewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, base.ErrorResponse(err))
		return
	}

	refreshPayload, err := auth.MainAuth.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, base.ErrorResponse(err))
		return
	}

	session, err := c.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if errors.Is(err, postgresdb.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, base.ErrorResponse(userRequests.ErrEmailOrPasswordNotCorrect))
			return
		}
		ctx.JSON(http.StatusInternalServerError, base.ErrorResponse(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusUnauthorized, base.ErrorResponse(err))
		return
	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, base.ErrorResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatched refresh token")
		ctx.JSON(http.StatusUnauthorized, base.ErrorResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		ctx.JSON(http.StatusUnauthorized, base.ErrorResponse(err))
		return
	}

	// create accessToken
	accessToken, accessPayload, err := auth.MainAuth.CreateToken(refreshPayload.Username, config.MainConfig.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, base.ErrorResponse(userRequests.ErrEmailOrPasswordNotCorrect))
		return
	}

	resp := userRequests.RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, resp)
}
