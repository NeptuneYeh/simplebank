package controllers

import (
	"github.com/NeptuneYeh/simplebank/init/logger"
	"github.com/NeptuneYeh/simplebank/init/store"
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
		store: *store.MainStore,
	}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var req userRequests.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	hashedPassword, err := hashPassword.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
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
				ctx.JSON(http.StatusForbidden, ErrorResponse(err))
				return
			}
		}
		logger.MainLog.Error().Msg(err.Error())
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	resp := userRequests.CreateUserResponse{
		Username:         user.Username,
		Email:            user.Email,
		FullName:         user.FullName,
		PasswordChangeAt: user.PasswordChangedAt,
		CreatedAt:        user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, resp)
}
