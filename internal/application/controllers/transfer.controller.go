package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/NeptuneYeh/simplebank/init/store"
	"github.com/NeptuneYeh/simplebank/internal/application/requests/transferRequests"
	postgresdb "github.com/NeptuneYeh/simplebank/internal/infrastructure/database/postgres/sqlc"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
)

type TransferController struct {
	store postgresdb.Store
}

func NewTransferController() *TransferController {
	return &TransferController{
		store: *store.MainStore,
	}
}

func (c *TransferController) CreateTransfer(ctx *gin.Context) {
	var req transferRequests.TransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !c.validAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}

	if !c.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	arg := postgresdb.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := c.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *TransferController) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := c.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}
