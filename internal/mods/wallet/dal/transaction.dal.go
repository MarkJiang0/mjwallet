package dal

import (
	"context"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/schema"
	"github.com/markjiang0/mjwallet/pkg/errors"
	"github.com/markjiang0/mjwallet/pkg/util"
	"gorm.io/gorm"
)

type Transaction struct {
	DB *gorm.DB
}

func GetTransactionDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(schema.Transaction))
}

func (t Transaction) Create(ctx context.Context, trans *schema.Transaction) error {
	res := GetTransactionDB(ctx, t.DB).Create(trans)
	return errors.WithStack(res.Error)
}

func (t Transaction) Update(ctx context.Context, item *schema.Transaction, selectFields ...string) error {
	db := GetTransactionDB(ctx, t.DB).Where("id=?", item.ID)
	if len(selectFields) > 0 {
		db = db.Select(selectFields)
	} else {
		db = db.Select("*").Omit("created_at")
	}
	result := db.Updates(item)
	return errors.WithStack(result.Error)
}

func (t Transaction) Query(ctx context.Context, params *schema.QueryTransactionParam, opts ...schema.TransactionQueryOptions) (*schema.TransactionQueryResult, error) {
	var opt schema.TransactionQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetTransactionDB(ctx, t.DB)
	if v := params.Status; v > 0 {
		db = db.Where("status = ?", v)
	}
	if v := params.FromAddress; len(v) > 0 {
		db = db.Where("from_address = ?", v)
	}
	if v := params.ToAddress; len(v) > 0 {
		db = db.Where("to_address = ?", v)
	}
	if v := params.TxHash; len(v) > 0 {
		db = db.Where("tx_hash = ?", v)
	}
	if v := params.Coin; len(v) > 0 {
		db = db.Where("coin = ?", v)
	}
	if v := params.UserId; len(v) > 0 {
		db = db.Where("user_id = ?", v)
	}

	var list schema.Transactions
	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := &schema.TransactionQueryResult{
		Data:       list,
		PageResult: pageResult,
	}
	return result, nil
}
