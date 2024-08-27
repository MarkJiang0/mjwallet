package dal

import (
	"context"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/schema"
	"github.com/markjiang0/mjwallet/pkg/errors"
	"github.com/markjiang0/mjwallet/pkg/util"
	"gorm.io/gorm"
	"strconv"
)

type Account struct {
	DB *gorm.DB
}

// Get user storage instance
func GetAccountDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(schema.Account))
}

func (a *Account) Create(ctx context.Context, item *schema.Account) error {
	res := GetAccountDB(ctx, a.DB).Create(item)
	return errors.WithStack(res.Error)
}

func (a Account) FindByUserID(ctx context.Context, userID string, opts ...util.QueryOptions) (*schema.Account, error) {
	var opt util.QueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	item := new(schema.Account)
	ok, err := util.FindOne(ctx, GetAccountDB(ctx, a.DB).Where("user_id=?", userID), opt, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

func (a Account) FindByAddressBase58(ctx context.Context, address string, opts ...util.QueryOptions) (*schema.Account, error) {
	var opt util.QueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	item := new(schema.Account)
	ok, err := util.FindOne(ctx, GetAccountDB(ctx, a.DB).Where("address_base58=?", address), opt, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

func (a Account) Query(ctx context.Context, param schema.QueryAccountParam, opts ...schema.AccountQueryOptions) (*schema.AccountQueryResult, error) {
	var opt schema.AccountQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	var list schema.AccountDtos

	db := GetAccountDB(ctx, a.DB).Select("account.id", "user.username", "account.address_base58",
		"account.balance", "account.usdt_balance", "account.created_at", "account.name", "account.default").
		Joins("join user on user.id = account.user_id")
	if v := param.LikeUserName; len(v) > 0 {
		db = db.Where("user.username LIKE ?", "%"+v+"%")
	}
	if v := param.UserId; len(v) > 0 {
		db = db.Where("account.user_id = ?", v)
	}
	if v := param.Default; len(v) > 0 {
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		db = db.Where("account.default = ?", val)
	}
	page, err := util.WrapPageQuery(ctx, db, param.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := &schema.AccountQueryResult{
		PageResult: page,
		Data:       list,
	}
	return result, nil
}

func (a Account) Update(ctx context.Context, item *schema.Account, selectFields ...string) error {
	db := GetAccountDB(ctx, a.DB).Where("id=?", item.ID)
	if len(selectFields) > 0 {
		db = db.Select(selectFields)
	} else {
		db = db.Select("*").Omit("created_at")
	}
	result := db.Updates(item)
	return errors.WithStack(result.Error)
}
