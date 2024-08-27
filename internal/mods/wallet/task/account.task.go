package task

import (
	"context"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/biz"
)

type Account struct {
	AccountBIZ *biz.Account
}

func (a Account) RefreshAccountBalance(ctx context.Context) error {
	err := a.AccountBIZ.RefreshAccountBalance(ctx)
	return err
}
