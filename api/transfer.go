package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/Just-A-NoobieDev/bankapi-gin-sqlc/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

// CreateTransfer godoc
//	@Summary		Create a new transfer
//	@Description	Create a new transfer between two accounts
//	@Param			transfer	body	createTransferRequest	true	"Create Transfer Request"
//	@Produce		application/json
//	@Tags			transfers
//	@Success		200	{object}	db.Transfer
//	@Router			/transfers [post]
func (server *Server) CreateTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}

	if !server.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}


	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transfer, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}

type getTransfersByAccountRequest struct {
	Id int64 `form:"id" binding:"required,min=1"`
	Page int32 `form:"page" binding:"required,min=1"`
	Size int32 `form:"size" binding:"required,min=1,max=10"`
}


// GetTransfersByAccount godoc
//	@Summary		Get transfers by account ID
//	@Description	Get transfers by the specified account ID
//	@Param			transfer	query	getTransfersByAccountRequest	true	"Get Transfers By Account Request"
//	@Produce		application/json
//	@Tags			transfers
//	@Success		200	{object}	[]db.Transfer
//	@Router			/transfers [get]
func (server *Server) GetTransfersByAccount(ctx *gin.Context) {
	var req getTransfersByAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	
	arg := db.GetTransfersByAccountParams{
		ID: req.Id,
		Off: (req.Page - 1) * req.Page,
		Size: req.Size,
	}

	transfers, err := server.store.GetTransfersByAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfers)
}

type getTransferByIdRequest struct {
	Id int64 `uri:"id" binding:"required,min=1"`
}

// GetTransferById godoc
//	@Summary		Get a transfer by ID
//	@Description	Get a transfer by the specified ID
//	@Param			id	path	getTransferByIdRequest	true	"Transfer ID"
//	@Produce		application/json
//	@Tags			transfers
//	@Success		200	{object}	db.Transfer
//	@Router			/transfers/{id} [get]
func (server *Server) GetTransferById(ctx *gin.Context) {
	var req getTransferByIdRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transfer, err := server.store.GetTransfer(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}


func (server *Server) validAccount(ctx *gin.Context, accountId int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account %d has different currency", accountId)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}