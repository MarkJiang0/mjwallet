package dal

import (
	"context"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/schema"
	"github.com/markjiang0/mjwallet/pkg/errors"
	"github.com/markjiang0/mjwallet/pkg/util"
	"gorm.io/gorm"
	"strconv"
)

type TelegramBot struct {
	DB *gorm.DB
}

func GetTelegramBotDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(schema.TelegramBot))
}

func (tb *TelegramBot) Create(ctx context.Context, item *schema.TelegramBot) error {
	res := GetTelegramBotDB(ctx, tb.DB).Create(item)
	return errors.WithStack(res.Error)
}

func (tb *TelegramBot) List(ctx context.Context, params *schema.QueryTelegramBotParam, opts ...*schema.TelegramBotQueryOpts) ([]*schema.TelegramBot, error) {
	var opt *schema.TelegramBotQueryOpts
	if len(opts) > 0 {
		opt = opts[0]
	}
	db := GetTelegramBotDB(ctx, tb.DB)
	if len(params.Name) > 0 {
		db = db.Where("name = ?", params.Name)
	}
	if len(params.Status) > 0 {
		v, err := strconv.Atoi(params.Status)
		if err != nil {
			return nil, err
		}
		db = db.Where("status = ?", v)
	}
	var list schema.TelegramBots
	_, err := util.FindList(ctx, db, opt.QueryOptions, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (tb *TelegramBot) Query(ctx context.Context, params *schema.QueryTelegramBotParam, opts ...*schema.TelegramBotQueryOpts) (*schema.QueryTelegramBotResult, error) {
	var opt *schema.TelegramBotQueryOpts
	if len(opts) > 0 {
		opt = opts[0]
	}
	db := GetTelegramBotDB(ctx, tb.DB)
	if len(params.Name) > 0 {
		db = db.Where("name = ?", params.Name)
	}
	if len(params.Status) > 0 {
		v, err := strconv.Atoi(params.Status)
		if err != nil {
			return nil, err
		}
		db = db.Where("status = ?", v)
	}
	var list schema.TelegramBots
	res, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, err
	}
	result := &schema.QueryTelegramBotResult{
		Data:       list,
		PageResult: res,
	}
	return result, nil
}

func (a *TelegramBot) Update(ctx context.Context, item *schema.TelegramBot, selectFields ...string) error {
	db := GetTelegramBotDB(ctx, a.DB).Where("id=?", item.ID)
	if len(selectFields) > 0 {
		db = db.Select(selectFields)
	} else {
		db = db.Select("*").Omit("created_at")
	}
	result := db.Updates(item)
	return errors.WithStack(result.Error)
}
