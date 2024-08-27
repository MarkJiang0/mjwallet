package schema

import (
	"github.com/markjiang0/mjwallet/internal/config"
	"github.com/markjiang0/mjwallet/pkg/util"
	"time"
)

type TelegramBot struct {
	ID        string    `json:"id" gorm:"size:20;primarykey;"`
	Name      string    `json:"name" gorm:"size:100;"`
	Token     string    `json:"token" gorm:"size:100"`
	Status    int8      `json:"status"`
	CreatedAt time.Time `json:"created_at" gorm:"index;"` // Create time
	UpdatedAt time.Time `json:"updated_at" gorm:"index;"` // Update time
}

func (a *TelegramBot) TableName() string {
	return config.C.FormatTableName("telegram_bot")
}

type QueryTelegramBotParam struct {
	util.PaginationParam
	Name   string `json:"name"`
	Status string `json:"status"`
}

type TelegramBotQueryOpts struct {
	util.QueryOptions
}

type TelegramBots []*TelegramBot

type QueryTelegramBotResult struct {
	Data       TelegramBots
	PageResult *util.PaginationResult
}
