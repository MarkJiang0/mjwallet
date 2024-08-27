package wallet

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/markjiang0/mjwallet/internal/config"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/api"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/biz"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/schema"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/task"
	"github.com/markjiang0/mjwallet/pkg/logging"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type WALLET struct {
	DB              *gorm.DB
	C               *cron.Cron
	AccountAPI      *api.Account
	TransactionAPI  *api.Transaction
	AccountTask     *task.Account
	TransactionTask *task.Transaction
	TelegramAPI     *api.TelegramBot
	TelegramTask    *task.TelegramBot
	TelegramBIZ     *biz.TelegramBot
}

func (a *WALLET) AutoMigrate(ctx context.Context) error {
	return a.DB.AutoMigrate(new(schema.Account), new(schema.Transaction), new(schema.TelegramBot))
}

func (a *WALLET) Init(ctx context.Context) error {
	if config.C.Storage.DB.AutoMigrate {
		if err := a.AutoMigrate(ctx); err != nil {
			return err
		}
	}
	if err := a.TelegramBIZ.LoadTelegramBot(ctx); err != nil {
		return err
	}
	return nil
}

func (a *WALLET) RegisterV1Routers(ctx context.Context, v1 *gin.RouterGroup) error {
	account := v1.Group("accounts")
	{
		account.POST("", a.AccountAPI.Create)
		account.GET("/get-account-by-address", a.AccountAPI.GetAccountByAddress)
		account.GET("", a.AccountAPI.Query)
		account.GET("/current", a.AccountAPI.QueryCurrentUserAccounts)
		account.PUT("/change-default", a.AccountAPI.ChangeDefaultAcount)
		account.GET("/refresh-balance", a.AccountAPI.RefreshAccountBalance)
	}
	transaction := v1.Group("transactions")
	{
		transaction.POST("", a.TransactionAPI.Transaction)
		transaction.GET("", a.TransactionAPI.Query)
	}
	bot := v1.Group("bots")
	{
		bot.POST("", a.TelegramAPI.Create)
		bot.GET("", a.TelegramAPI.Query)
	}
	return nil
}

func (w *WALLET) RegisterAndStartCronJob(ctx context.Context) error {
	_, err := w.C.AddFunc("0 */1 * * *", func() {
		err := w.AccountTask.RefreshAccountBalance(ctx)
		if err != nil {
			logging.Context(ctx).Error("failed to add task", zap.Error(err))
		}
	})
	if err != nil {
		return err
	}

	_, err = w.C.AddFunc("0/1 * * * *", func() {
		err := w.TransactionTask.ConfirmTransaction(ctx)
		if err != nil {
			logging.Context(ctx).Error("failed to add task", zap.Error(err))
		}
	})
	if err != nil {
		return err
	}

	_, err = w.C.AddFunc("0/10 * * * *", func() {
		err := w.TelegramTask.StartNewBot(ctx)
		if err != nil {
			logging.Context(ctx).Error("failed to add task", zap.Error(err))
		}
	})
	if err != nil {
		return err
	}

	w.C.Start()
	return nil
}

func (a *WALLET) Release(ctx context.Context) error {
	return nil
}
