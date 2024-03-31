package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/Just-A-NoobieDev/bankapi-gin-sqlc/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Name     string `json:"name" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

// CreateAccount		godoc
//	@Summary		Create a new account
//	@Description	Create a new account with the specified name and currency
//	@Param			account	body	createAccountRequest	true	"Create Account Request"
//	@Produce		application/json
//	@Tags			accounts
//	@Success		200	{object}	db.Account
//	@Router			/acounts [post]
func (server *Server) CreateAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Name:    req.Name,
		Currency: req.Currency,
		Balance: 0,
	}
	
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// GetAccount		godoc
//	@Summary		Get an account by ID
//	@Description	Get an account by the specified ID
//	@Param			id	path	getAccountRequest	true	"Account ID"
//	@Produce		application/json
//	@Tags			accounts
//	@Success		200	{object}	db.Account
//	@Router			/accounts/{id} [get]
func (server *Server) GetAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountsParams struct {
	Page  int32 `form:"page" binding:"required,min=1"`
	Size int32 `form:"size" binding:"required,min=5,max=10"`
}

// GetAccounts		godoc
//	@Summary		Get a list of accounts
//	@Description	Get a list of accounts with pagination
//	@Param			pagination	query	getAccountsParams	true	"Pagination"
//	@Produce		application/json
//	@Tags			accounts
//	@Success		200	{object}	[]db.Account
//	@Router			/accounts [get]
func (server *Server) GetAccounts(ctx *gin.Context) {

	var req getAccountsParams
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetAccountsParams{
		Limit:  req.Size,
		Offset: (req.Page - 1) * req.Page,
	}

	accounts, err := server.store.GetAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	fmt.Println(accounts)

	ctx.JSON(http.StatusOK, accounts)
}

type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}


// DeleteAccount		godoc
//	@Summary		Delete an account by ID
//	@Description	Delete an account by the specified ID
//	@Param			id	path	deleteAccountRequest	true	"Account ID"
//	@Produce		application/json
//	@Tags			accounts
//	@Success		200	{object}	string
//	@Router			/accounts/{id} [delete]
func (server *Server) DeleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}



	ctx.JSON(http.StatusOK, gin.H{"success": "Account deleted successfully!"})
}

type depositRequest struct {
	ID int64 `json:"id" binding:"required,min=1"`
	Amount int64 `json:"amount" binding:"required,gt=0"`
}


// Deposit		godoc
//	@Summary		Deposit money to an account
//	@Description	Deposit money to an account by the specified ID
//	@Param			account	body	depositRequest	true	"Deposit Request"
//	@Produce		application/json
//	@Tags			accounts
//	@Success		200	{object}	db.Account
//	@Router			/accounts/deposit [post]
func (server *Server) Deposit(ctx *gin.Context) {
	var req depositRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.AddAccountBalanceParams{
		ID: req.ID,
		Amount: req.Amount,
	}

	account, err := server.store.AddAccountBalance(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}