package task

import (
	"context"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/biz"
)

type TelegramBot struct {
	TelegramBotBIZ *biz.TelegramBot
}

func (tb *TelegramBot) StartNewBot(ctx context.Context) error {
	return tb.TelegramBotBIZ.StartNewBots(ctx)
}
