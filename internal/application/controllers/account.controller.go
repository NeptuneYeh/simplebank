package controllers

import (
	"errors"
	"github.com/NeptuneYeh/simplebank/init/store"
	"github.com/NeptuneYeh/simplebank/internal/application/base"
	"github.com/NeptuneYeh/simplebank/internal/application/middlewares"
	"github.com/NeptuneYeh/simplebank/internal/application/requests/accountRequests"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/NeptuneYeh/simplebank/tools/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AccountController struct {
	store postgresdb.Store
}

func NewAccountController() *AccountController {
	return &AccountController{
		store: store.MainStore,
	}
}

func (c *AccountController) CreateAccount(ctx *gin.Context) {
	var req accountRequests.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, base.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

	arg := postgresdb.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := c.store.CreateAccount(ctx, arg)
	if err != nil {
		switch postgresdb.ErrorCode(err) {
		case postgresdb.ForeignKeyViolation:
		case postgresdb.UniqueViolation:
			ctx.JSON(http.StatusForbidden, base.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, base.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (c *AccountController) GetAccount(ctx *gin.Context) {
	var req accountRequests.GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, base.ErrorResponse(err))
		return
	}

	account, err := c.store.GetAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, postgresdb.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, base.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, base.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, base.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (c *AccountController) ListAccount(ctx *gin.Context) {
	var req accountRequests.ListAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, base.ErrorResponse(err))
		return
	}

	authPayload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
	arg := postgresdb.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	account, err := c.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, base.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
