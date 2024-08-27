package api

import (
	"github.com/gin-gonic/gin"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/biz"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/schema"
	"github.com/markjiang0/mjwallet/pkg/util"
)

type TelegramBot struct {
	TelegramBotBIZ *biz.TelegramBot
}

// @Tags TelegramBotAPI
// @Security ApiKeyAuth
// @Summary Create TelegramBot record
// @Success 200 {object} util.ResponseResult{data=schema.Account}
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/bots [post]
func (a *TelegramBot) Create(c *gin.Context) {
	var bot schema.TelegramBot
	err := util.ParseJSON(c, &bot)
	if err != nil {
		util.ResError(c, err)
		return
	}

	err = a.TelegramBotBIZ.Create(c.Request.Context(), &bot)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResSuccess(c, nil)
}

// @Tags TelegramBotAPI
// @Security ApiKeyAuth
// @Summary Query bot list
// @Success 200 {object} util.ResponseResult{data=schema.Account}
// @Failure 400 {object} util.ResponseResult
// @Failure 401 {object} util.ResponseResult
// @Failure 500 {object} util.ResponseResult
// @Router /api/v1/bots [get]
func (a *TelegramBot) Query(c *gin.Context) {
	var param schema.QueryTelegramBotParam
	if err := util.ParseQuery(c, &param); err != nil {
		util.ResError(c, err)
		return
	}
	result, err := a.TelegramBotBIZ.Query(c.Request.Context(), &param)
	if err != nil {
		util.ResError(c, err)
		return
	}
	util.ResPage(c, result.Data, result.PageResult)
}
