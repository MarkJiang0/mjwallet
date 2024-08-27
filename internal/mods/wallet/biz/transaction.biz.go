package biz

import (
	"context"
	"fmt"
	tronWallet "github.com/criptoDevs/tron-wallet"
	"github.com/criptoDevs/tron-wallet/enums"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/dal"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/schema"
	"github.com/markjiang0/mjwallet/pkg/logging"
	"github.com/markjiang0/mjwallet/pkg/util"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Transaction struct {
	TransactionDAL *dal.Transaction
	AccountDAL     *dal.Account
}

func (t *Transaction) DoTransaction(ctx context.Context, transaction *schema.Transaction, checkOwner bool) error {
	account, err := t.AccountDAL.FindByAddressBase58(ctx, transaction.FromAddress)
	if err != nil {
		return err
	}
	if checkOwner {
		userID := util.FromUserID(ctx)
		if userID != account.UserId {
			return fmt.Errorf("地址不属于当前用户")
		}
	}

	wallet, err := tronWallet.CreateTronWallet(enums.MAIN_NODE, account.PrivateKey)
	if err != nil {
		return err
	}
	trxBalance, err := wallet.Balance()
	if err != nil {
		return err
	}

	if transaction.Coin == "TRX" {
		fee, err := wallet.EstimateTransferFee(transaction.ToAddress, transaction.Amount)
		if err != nil {
			return err
		}
		if trxBalance <= (transaction.Amount + fee) {
			return fmt.Errorf("钱包TRX余额不足")
		}
		txHash, err := wallet.Transfer(transaction.ToAddress, transaction.Amount)
		if err != nil {
			return err
		}
		transaction.TxHash = txHash
	} else if transaction.Coin == "USDT" {
		token := &tronWallet.Token{
			ContractAddress: enums.MAIN_Tether_USDT,
		}
		balance, err := wallet.BalanceTRC20(token)
		if err != nil {
			return err
		}
		fee, err := wallet.EstimateTransferTRC20Fee()
		if err != nil {
			return err
		}
		if balance <= transaction.Amount {
			return fmt.Errorf("钱包USDT余额不足")
		}
		if trxBalance <= fee {
			return fmt.Errorf("钱包TRX余额不足")
		}
		txHash, err := wallet.TransferTRC20(token, transaction.ToAddress, transaction.Amount)
		transaction.TxHash = txHash
	}

	transaction.ID = util.NewXID()
	transaction.UserId = account.UserId
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()
	transaction.Status = 1
	transaction.ConfirmStatus = 0
	err = t.TransactionDAL.Create(ctx, transaction)
	if err != nil {
		return err
	}
	return nil
}

func (t Transaction) ConfirmTransaction(ctx context.Context) error {
	params := &schema.QueryTransactionParam{
		Status: 0,
	}
	queryResult, err := t.TransactionDAL.Query(ctx, params)
	if err != nil {
		return err
	}
	txMap := make(map[string][]*schema.Transaction)
	var addresses []string
	if len(queryResult.Data) > 0 {
		for _, tx := range queryResult.Data {
			if _, ok := txMap[tx.FromAddress]; ok {
				txMap[tx.FromAddress] = append(txMap[tx.FromAddress], tx)
			} else {
				txMap[tx.FromAddress] = []*schema.Transaction{tx}
			}

			if addresses != nil && len(addresses) > 0 {
				exist := false
				for _, a := range addresses {
					if a == tx.FromAddress {
						exist = true
						break
					}
				}
				if !exist {
					addresses = append(addresses, tx.FromAddress)
				}
			} else {
				addresses = []string{tx.FromAddress}
			}
		}
	}

	err = t.DoConfirmTransaction(ctx, addresses, txMap)
	if err != nil {
		return err
	}
	return nil
}

func (t *Transaction) DoConfirmTransaction(ctx context.Context, addresses []string, txMap map[string][]*schema.Transaction) error {
	c := &tronWallet.Crawler{
		Node:      enums.MAIN_NODE,
		Addresses: addresses,
	}
	res, err := c.ScanBlocks(10)
	if err != nil {
		return err
	}

	if len(res) > 0 {
		var wg sync.WaitGroup
		for _, txs := range txMap {
			wg.Add(1)
			go func(transList []*schema.Transaction) {
				defer wg.Done()

				for _, r := range res {
					for _, ct := range r.Transactions {
						for _, tx := range transList {
							if ct.TxId == tx.TxHash {
								if ct.Confirmations > 3 {
									logging.Context(ctx).Info("Confirmed transaction", zap.String("txhash", tx.TxHash))
									tx.ConfirmStatus = 1
									if err := t.TransactionDAL.Update(ctx, tx, "status"); err != nil {
										fmt.Println(err)
										return
									}
								}
							}
						}
					}
				}

			}(txs)
		}
		wg.Wait()

	}
	return nil
}

func (a Transaction) Query(ctx context.Context, param schema.QueryTransactionParam) (*schema.TransactionQueryResult, error) {
	query, err := a.TransactionDAL.Query(ctx, &param, schema.TransactionQueryOptions{
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
