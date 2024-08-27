package api

import (
	"github.com/gin-gonic/gin"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/biz"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/schema"
	"github.com/markjiang0/mjwallet/pkg/util"
)

type Transaction struct {
	TransactionBIZ *biz.Transaction
}

// @Tags TransactionAPI
// @Security ApiKeyAuth
// @Summary Create transacction record
// @Success 200 {object} util.ResponseResult{data=schema.Account}
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/transactions [post]
func (a *Transaction) Transaction(c *gin.Context) {
	var param schema.Transaction
	if err := util.ParseJSON(c, &param); err != nil {
		util.ResError(c, err)
		return
	}
	err := a.TransactionBIZ.DoTransaction(c.Request.Context(), &param, true)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, nil)
}

// @Tags TransactionAPI
// @Security ApiKeyAuth
// @Summary Query transaction list
// @Success 200 {object} util.ResponseResult{data=schema.Account}
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/transactions [get]
func (a *Transaction) Query(c *gin.Context) {
	var param schema.QueryTransactionParam
	if err := util.ParseQuery(c, &param); err != nil {
		util.ResError(c, err)
		return
	}
	result, err := a.TransactionBIZ.Query(c.Request.Context(), param)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResPage(c, result.Data, result.PageResult)
}
