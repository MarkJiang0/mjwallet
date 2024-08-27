package biz

import (
	"context"
	"fmt"
	tronWallet "github.com/criptoDevs/tron-wallet"
	"github.com/criptoDevs/tron-wallet/enums"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/dal"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/schema"
	"github.com/markjiang0/mjwallet/pkg/cachex"
	"github.com/markjiang0/mjwallet/pkg/util"
	"time"
)

type Account struct {
	Cache      cachex.Cacher
	Trans      *util.Trans
	AccountDAL *dal.Account
}

func (a *Account) Create(ctx context.Context, name string) (*schema.Account, error) {
	userID := util.FromUserID(ctx)
	wallet := tronWallet.GenerateTronWallet(enums.MAIN_NODE)
	id := util.NewXID()
	if name == "" {
		name = "account " + id
	}
	account := &schema.Account{
		ID:            id,
		UserId:        userID,
		Address:       wallet.Address,
		AddressBase58: wallet.AddressBase58,
		PublicKey:     wallet.PublicKey,
		PrivateKey:    wallet.PrivateKey,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Default:       -1,
		Name:          name,
	}
	err := a.AccountDAL.Create(ctx, account)

	if err != nil {
		return nil, err
	}
	return account, nil
}

func (a *Account) GetAccountBalance(account *schema.Account) (int64, int64, error) {

	wallet, err := tronWallet.CreateTronWallet(enums.MAIN_NODE, account.PrivateKey)
	if err != nil {
		return 0, 0, err
	}
	balance, err := wallet.Balance()

	//c, err := grpcClient.GetGrpcClient(enums.MAIN_NODE)
	//acc, err := c.GetAccount("TUzCLs7Ab41QRhcYy5mj8jsLLonYkdNcoh")
	if err != nil {
		return balance, 0, err
	}
	token := &tronWallet.Token{
		ContractAddress: enums.MAIN_Tether_USDT,
	}
	trc20Balance, err := wallet.BalanceTRC20(token)
	if err != nil {
		return balance, trc20Balance, err
	}
	return balance, trc20Balance, nil
}

func (a *Account) GetAccountByAddress(ctx context.Context, address string) (*schema.AccountDto, error) {
	userID := util.FromUserID(ctx)
	account, err := a.AccountDAL.FindByAddressBase58(ctx, address)
	if userID != account.UserId {
		return nil, fmt.Errorf("地址不属于当前用户")
	}
	if err != nil {
		return nil, err
	}
	balance, usdtBalance, err := a.GetAccountBalance(account)
	if err != nil {
		return nil, err
	}
	accountDto := &schema.AccountDto{
		ID:            account.ID,
		AddressBase58: account.AddressBase58,
		Balance:       balance,
		UsdtBalance:   usdtBalance,
		CreatedAt:     account.CreatedAt,
		Name:          account.Name,
	}

	account.Balance = balance
	account.UsdtBalance = usdtBalance
	err = a.AccountDAL.Update(ctx, account, "balance", "usdt_balance")
	if err != nil {
		return nil, err
	}

	return accountDto, nil
}

func (a *Account) transfer(ctx context.Context, param *schema.TransferParam) (string, error) {
	userID := util.FromUserID(ctx)
	account, err := a.AccountDAL.FindByUserID(ctx, userID)
	if err != nil {
		return "", err
	}
	wallet, err := tronWallet.CreateTronWallet(enums.MAIN_NODE, account.PrivateKey)
	if err != nil {
		return "", err
	}
	txId, err := wallet.Transfer(param.ToAddress, param.Amount)

	if err != nil {
		return "", err
	}
	return txId, err
}

func (a *Account) estimateTransferFee(ctx context.Context, param *schema.TransferParam) (int64, error) {
	userID := util.FromUserID(ctx)
	account, err := a.AccountDAL.FindByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	wallet, err := tronWallet.CreateTronWallet(enums.MAIN_NODE, account.PrivateKey)
	if err != nil {
		return 0, err
	}
	fee, err := wallet.EstimateTransferFee(param.ToAddress, param.Amount)
	if err != nil {
		return 0, err
	}
	return fee, err
}

func (a Account) Query(ctx context.Context, param schema.QueryAccountParam) (*schema.AccountQueryResult, error) {
	query, err := a.AccountDAL.Query(ctx, param, schema.AccountQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: []util.OrderByParam{
				{Field: "account.created_at", Direction: util.DESC},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return query, nil
}

func (a Account) ChangeDefaultAcount(ctx context.Context, address string) error {
	userID := util.FromUserID(ctx)
	account, err := a.AccountDAL.FindByAddressBase58(ctx, address)
	if userID != account.UserId {
		return fmt.Errorf("地址不属于当前用户")
	}
	if err != nil {
		return err
	}

	param := schema.QueryAccountParam{
		Default: "1",
	}
	query, err := a.AccountDAL.Query(ctx, param, schema.AccountQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: []util.OrderByParam{
				{Field: "account.created_at", Direction: util.DESC},
			},
		},
	})
	if err != nil {
		return err
	}

	//a.AccountDAL.DB.Begin()

	for _, account := range query.Data {
		if account.Default {
			item := &schema.Account{
				ID:      account.ID,
				Default: 0,
			}
			err := a.AccountDAL.Update(ctx, item, "Default")
			if err != nil {
				return err
			}
		}
	}

	account.Default = 1
	err = a.AccountDAL.Update(ctx, account, "Default")
	if err != nil {
		//a.AccountDAL.DB.Rollback()
		return err
	}

	//a.AccountDAL.DB.Commit()
	return nil
}

func (a *Account) RefreshAccountBalance(ctx context.Context) error {
	paginationParam := util.PaginationParam{
		PageSize: 9999,
	}
	param := schema.QueryAccountParam{
		PaginationParam: paginationParam,
	}
	query, err := a.AccountDAL.Query(ctx, param, schema.AccountQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: []util.OrderByParam{
				{Field: "account.created_at", Direction: util.DESC},
			},
		},
	})
	if err != nil {
		return err
	}
	for _, accountDto := range query.Data {
		account, err := a.AccountDAL.FindByAddressBase58(ctx, accountDto.AddressBase58)
		if err != nil {
			return err
		}
		balance, usdtBalance, err := a.GetAccountBalance(account)
		account.Balance = balance
		account.UsdtBalance = usdtBalance
		err = a.AccountDAL.Update(ctx, account, "balance", "usdt_balance")
	}
	return nil
}
