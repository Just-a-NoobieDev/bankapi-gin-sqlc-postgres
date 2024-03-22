package api

import (
	"net/http"

	db "github.com/Just-A-NoobieDev/bankapi-gin-sqlc/db/sqlc"
	"github.com/gin-gonic/gin"
)

type getEntriesByAccountRequest struct {
	Id   int64 `form:"id" binding:"required,min=1"`
	Page int32 `form:"page" binding:"required,min=1"`
	Size int32 `form:"size" binding:"required,min=1,max=10"`
}

func (server *Server) GetEntriesByAccount(ctx *gin.Context) {
	var req getEntriesByAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetEntriesParams{
		AccountID: req.Id,
		Limit:     req.Size,
		Offset:    (req.Page - 1) * req.Page,
	}

	entries, err := server.store.GetEntries(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)
}