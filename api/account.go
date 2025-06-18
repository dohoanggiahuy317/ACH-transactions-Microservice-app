package api

import (
	"database/sql"
	"net/http"

	db "github.com/dohoanggiahuy317/ACH-transactions-Microservice-app/db/sqlc"
	"github.com/gin-gonic/gin"
)

// ======== Create Account Request Struct ========
// createAccountRequest defines the structure for creating a new account.
type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createAccount(ctx *gin.Context) {

	// Create a new account based on the request body.
	var req createAccountRequest

	// Bind the JSON request body to the createAccountRequest struct.
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Create the account using the store's CreateAccount method.

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Respond with the created account in JSON format.
	ctx.JSON(http.StatusOK, account)

}


// ========= Get Account Request Struct ========
// getAccountRequest defines the structure for retrieving an account by ID.

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest

	// Bind the URI parameter to the getAccountRequest struct.
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Retrieve the account using the store's GetAccount method.
	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		// If an unexpected error occurs, respond with a 500 Internal Server Error.
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Respond with the retrieved account in JSON format.
	ctx.JSON(http.StatusOK, account)
}


// ========== List Accounts Request Struct =========
// listAccountsRequest defines the structure for listing accounts with pagination.

type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest

	// Bind the query parameters to the listAccountsRequest struct.
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// List accounts using the store's ListAccounts method.
	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Respond with the list of accounts in JSON format.
	ctx.JSON(http.StatusOK, accounts)
}