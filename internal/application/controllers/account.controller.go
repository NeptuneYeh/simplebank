package controllers

import (
	"database/sql"
	"errors"
	"github.com/NeptuneYeh/simplebank/init/store"
	"github.com/NeptuneYeh/simplebank/internal/application/requests/accountRequests"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"net/http"
)

type AccountController struct {
	store postgresdb.Store
}

func NewAccountController() *AccountController {
	return &AccountController{
		store: *store.MainStore,
	}
}

func (c *AccountController) CreateAccount(ctx *gin.Context) {
	var req accountRequests.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	arg := postgresdb.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := c.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			//logger.MainLog.Error().Msg(pqErr.Code.Name())
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, ErrorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (c *AccountController) GetAccount(ctx *gin.Context) {
	var req accountRequests.GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	account, err := c.store.GetAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (c *AccountController) ListAccount(ctx *gin.Context) {
	var req accountRequests.ListAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	arg := postgresdb.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	account, err := c.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
