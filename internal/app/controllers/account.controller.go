package controllers

import (
	"database/sql"
	"github.com/NeptuneYeh/simplerecommend/init/asynq"
	"github.com/NeptuneYeh/simplerecommend/init/auth"
	"github.com/NeptuneYeh/simplerecommend/init/config"
	"github.com/NeptuneYeh/simplerecommend/init/store"
	"github.com/NeptuneYeh/simplerecommend/internal/app/requests/accountRequests"
	db "github.com/NeptuneYeh/simplerecommend/internal/infra/database/mysql/sqlc"
	"github.com/NeptuneYeh/simplerecommend/pkg/hashPassword"
	"github.com/NeptuneYeh/simplerecommend/pkg/token"
	"github.com/NeptuneYeh/simplerecommend/worker"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type AccountController struct {
	store           db.Store
	auth            token.Maker
	config          *config.Module
	taskDistributor worker.TaskDistributor
}

func NewAccountController() *AccountController {
	return &AccountController{
		store:           *store.MyStore,
		auth:            *auth.MyAuth,
		config:          config.MyConfig,
		taskDistributor: asynq.MyAsynq.TaskDistributor,
	}
}

func (c *AccountController) CreateAccount(ctx *gin.Context) {
	var req accountRequests.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	hashedPassword, err := hashPassword.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Name:     req.Name,
		Password: hashedPassword,
		Email:    req.Email,
		IsValid:  false,
	}

	result, err := c.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	getAccount, err := c.store.GetAccount(ctx, lastID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	resp := accountRequests.AccountResponse{
		ID:        getAccount.ID,
		Email:     getAccount.Email,
		Name:      getAccount.Name,
		IsValid:   getAccount.IsValid,
		CreatedAt: getAccount.CreatedAt,
		UpdatedAt: getAccount.UpdatedAt,
	}

	// create email valid record
	code := uuid.New().String()
	argValid := db.CreateAccountValidParams{
		AccountID: getAccount.ID,
		Code:      code,
		ExpiredAt: time.Now().Add(15 * time.Minute),
	}

	_, err = c.store.CreateAccountValid(ctx, argValid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	// send email
	taskPayload := &worker.PayloadSendVerifyEmail{
		Name:  getAccount.Name,
		Email: getAccount.Email,
		Code:  code,
	}
	err = c.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload)
	if err != nil {
		// TODO logger 即可 驗證信寄不出去應該也不能妨礙 user 繼續登入，通常是登入之後提醒尚未驗證
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *AccountController) LoginAccount(ctx *gin.Context) {
	var req accountRequests.LoginAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	account, err := c.store.GetAccountByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ErrorResponse(accountRequests.ErrEmailOrPasswordNotCorrect))
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	// check password
	err = hashPassword.CheckPassword(req.Password, account.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(accountRequests.ErrEmailOrPasswordNotCorrect))
		return
	}

	// create accessToken
	accessToken, err := c.auth.CreateToken(account.Name, c.config.AccessTokenDuration)

	resp := accountRequests.LoginAccountResponse{
		AccessToken: accessToken,
		Email:       account.Email,
		Name:        account.Name,
		IsValid:     account.IsValid,
		CreatedAt:   account.CreatedAt,
		UpdatedAt:   account.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *AccountController) VerifyEmail(ctx *gin.Context) {
	var req accountRequests.VerifyEmailRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	arg := db.VerifyEmailTxParams{
		Code: req.Code,
	}

	txResult, err := c.store.VerifyEmailTx(ctx, arg)
	if err != nil {
		//ctx.JSON(http.StatusForbidden, ErrorResponse(accountRequests.ErrVerifyEmail))
		ctx.JSON(http.StatusForbidden, ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, txResult)
}
