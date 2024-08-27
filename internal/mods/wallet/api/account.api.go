package api

import (
	"github.com/gin-gonic/gin"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/biz"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/schema"
	"github.com/markjiang0/mjwallet/pkg/util"
)

type Account struct {
	AccountBIZ *biz.Account
}

// @Tags AccountAPI
// @Security ApiKeyAuth
// @Summary Create account record
// @Success 200 {object} util.ResponseResult{data=schema.Account}
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/accounts [post]
func (a *Account) Create(c *gin.Context) {
	name := c.Param("name")
	account, err := a.AccountBIZ.Create(c.Request.Context(), name)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, account)
}

// @Tags AccountAPI
// @Security ApiKeyAuth
// @Summary Get account
// @Success 200 {object} util.ResponseResult{data=schema.Account}
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/accounts/get-account-by-address [get]
func (a *Account) GetAccountByAddress(c *gin.Context) {
	address := c.Query("address")
	account, err := a.AccountBIZ.GetAccountByAddress(c.Request.Context(), address)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, account)
}

// @Tags AccountAPI
// @Security ApiKeyAuth
// @Summary Change default account
// @Success 200 {object} util.ResponseResult{data=schema.Account}
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/accounts/change-default [put]
func (a *Account) ChangeDefaultAcount(c *gin.Context) {
	address := c.Query("address")
	err := a.AccountBIZ.ChangeDefaultAcount(c.Request.Context(), address)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, nil)
}

// @Tags AccountAPI
// @Security ApiKeyAuth
// @Summary Query account list
// @Success 200 {object} util.ResponseResult{data=schema.Account}
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/accounts [get]
func (a *Account) Query(c *gin.Context) {
	var param schema.QueryAccountParam
	if err := util.ParseQuery(c, &param); err != nil {
		util.ResError(c, err)
		return
	}
	result, err := a.AccountBIZ.Query(c.Request.Context(), param)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResPage(c, result.Data, result.PageResult)
}

// @Tags AccountAPI
// @Security ApiKeyAuth
// @Summary Query current user account list
// @Success 200 {object} util.ResponseResult{data=schema.Account}
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/accounts/current [get]
func (a *Account) QueryCurrentUserAccounts(c *gin.Context) {
	userID := util.FromUserID(c.Request.Context())
	param := schema.QueryAccountParam{
		UserId: userID,
	}

	result, err := a.AccountBIZ.Query(c.Request.Context(), param)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResPage(c, result.Data, result.PageResult)
}

// @Tags AccountAPI
// @Security ApiKeyAuth
// @Summary Refresh Account balance
// @Success 200 {object} util.ResponseResult{data=schema.Account}
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/accounts/refresh-balance [get]
func (a Account) RefreshAccountBalance(c *gin.Context) {
	err := a.AccountBIZ.RefreshAccountBalance(c.Request.Context())
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, nil)
}

// 更新balance job
// 记录交易 查询交易
