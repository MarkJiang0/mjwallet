package biz

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/dal"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/schema"
	"github.com/markjiang0/mjwallet/pkg/encoding/json"
	"github.com/markjiang0/mjwallet/pkg/logging"
	"github.com/markjiang0/mjwallet/pkg/util"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type TelegramBot struct {
	TelegramBotDAL *dal.TelegramBot
	TransactionBIZ *Transaction
	AccountBIZ     *Account
}

func (tb *TelegramBot) Create(ctx context.Context, telegramBot *schema.TelegramBot) error {

	telegramBot.CreatedAt = time.Now()
	telegramBot.UpdatedAt = time.Now()
	telegramBot.ID = util.NewXID()
	telegramBot.Status = 0

	err := tb.TelegramBotDAL.Create(ctx, telegramBot)
	if err != nil {
		return err
	}

	//err = tb.startBot(ctx, telegramBot)
	return err
}

func (tb *TelegramBot) LoadTelegramBot(ctx context.Context) error {
	params := &schema.QueryTelegramBotParam{
		Status: "1",
	}
	list, err := tb.TelegramBotDAL.List(ctx, params, &schema.TelegramBotQueryOpts{
		QueryOptions: util.QueryOptions{
			OrderFields: []util.OrderByParam{
				{Field: "created_at", Direction: util.DESC},
			},
		},
	})
	if err != nil {
		return err
	}
	for _, telegramBot := range list {
		err = tb.startBot(ctx, telegramBot)
		if err != nil {
			return err
		}
	}

	//err = tb.startNewBots(ctx)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (tb *TelegramBot) StartNewBots(ctx context.Context) error {
	// 新建的bot
	params := &schema.QueryTelegramBotParam{
		Status: "0",
	}
	list, err := tb.TelegramBotDAL.List(ctx, params, &schema.TelegramBotQueryOpts{
		QueryOptions: util.QueryOptions{
			OrderFields: []util.OrderByParam{
				{Field: "created_at", Direction: util.DESC},
			},
		},
	})
	if err != nil {
		return err
	}
	for _, telegramBot := range list {
		err = tb.startBot(ctx, telegramBot)
		telegramBot.Status = 1
		err = tb.TelegramBotDAL.Update(ctx, telegramBot)
		if err != nil {
			return err
		}
	}
	return nil
}

func (tb *TelegramBot) startBot(ctx context.Context, telegramBot *schema.TelegramBot) error {
	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
		bot.WithMessageTextHandler("/transfer", bot.MatchTypePrefix, tb.HandlerTransferCommand),
	}
	b, err := bot.New(telegramBot.Token, opts...)
	if err != nil {
		return err
	}
	go func() {
		logging.Context(ctx).Info("Bot starting", zap.String("name", telegramBot.Name))
		b.Start(ctx)
	}()
	return nil
}

func (tb *TelegramBot) HandlerTransferCommand(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		logging.Context(ctx).Info(json.MarshalToString(update.Message))
		strArr := strings.Split(update.Message.Text, "\n")
		if len(strArr) == 3 {
			address := strArr[1]
			amount, err := strconv.Atoi(strArr[2])
			if err != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "invalid param",
				})
				return
			}
			queryRes, err := tb.AccountBIZ.Query(ctx, schema.QueryAccountParam{
				UserId: "cr3hvplupt9nb3n2bq40",
			})
			var fromAddress string
			if len(queryRes.Data) > 0 {
				for _, a := range queryRes.Data {
					if a.Default {
						fromAddress = a.AddressBase58
					}
				}
			}
			transaction := &schema.Transaction{
				ToAddress:   address,
				Amount:      int64(amount),
				Coin:        "TRX",
				FromAddress: fromAddress,
			}
			err = tb.TransactionBIZ.DoTransaction(ctx, transaction, false)
			if err != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "transfer failed",
				})
				return
			}
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "transfer success!",
			})
		}

	}
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "本机器人只支持转账操作,请使用一下格式进行转账:\n" + "/transfer\n" + "{trc20 address}\n" + "{amount}\n",
		})
	}
}

func (a *TelegramBot) Query(ctx context.Context, param *schema.QueryTelegramBotParam) (*schema.QueryTelegramBotResult, error) {
	query, err := a.TelegramBotDAL.Query(ctx, param, &schema.TelegramBotQueryOpts{
		QueryOptions: util.QueryOptions{
			OrderFields: []util.OrderByParam{
				{Field: "created_at", Direction: util.DESC},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return query, nil
}
