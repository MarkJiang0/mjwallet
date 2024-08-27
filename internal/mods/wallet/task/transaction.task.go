package task

import (
	"context"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/biz"
)

type Transaction struct {
	TransactionBIZ *biz.Transaction
}

func (t Transaction) ConfirmTransaction(ctx context.Context) error {
	return t.TransactionBIZ.ConfirmTransaction(ctx)
}
