package controllers

import (
	"errors"
	"fmt"
	"github.com/NeptuneYeh/simplebank/init/store"
	"github.com/NeptuneYeh/simplebank/internal/application/base"
	"github.com/NeptuneYeh/simplebank/internal/application/middlewares"
	"github.com/NeptuneYeh/simplebank/internal/application/requests/transferRequests"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/NeptuneYeh/simplebank/tools/token"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
)

type TransferController struct {
	store postgresdb.Store
}

func NewTransferController() *TransferController {
	return &TransferController{
		store: store.MainStore,
	}
}

func (c *TransferController) CreateTransfer(ctx *gin.Context) {
	var req transferRequests.TransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, base.ErrorResponse(err))
		return
	}

	fromAccount, valid := c.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, base.ErrorResponse(err))
		return
	}
	_, valid = c.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	arg := postgresdb.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := c.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, base.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *TransferController) validAccount(ctx *gin.Context, accountID int64, currency string) (postgresdb.Account, bool) {
	account, err := c.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, postgresdb.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, base.ErrorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, base.ErrorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, base.ErrorResponse(err))
		return account, false
	}

	return account, true
}
