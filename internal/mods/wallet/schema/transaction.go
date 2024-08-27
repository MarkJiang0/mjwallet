package schema

import (
	"github.com/markjiang0/mjwallet/internal/config"
	"github.com/markjiang0/mjwallet/internal/mods/rbac/schema"
	"github.com/markjiang0/mjwallet/pkg/util"
	"time"
)

type Transaction struct {
	ID            string    `json:"id" gorm:"size:20;primarykey;"`
	UserId        string    `json:"user_id" gorm:"size:20"`
	FromAddress   string    `json:"from_address" gorm:"size:64"`
	ToAddress     string    `json:"to_address" gorm:"size:64"`
	Coin          string    `json:"coin" gorm:"size:10"`
	Amount        int64     `json:"amount,string,omitempty"`
	TxHash        string    `json:"tx_hash" gorm:"size:128"`
	Status        int8      `json:"status"`                   // Status 0-pending 1-success 2-failed
	ConfirmStatus int8      `json:"confirm_status"`           // Status 0-pending 1-confirmed
	CreatedAt     time.Time `json:"created_at" gorm:"index;"` // Create time
	UpdatedAt     time.Time `json:"updated_at" gorm:"index;"` // Update time
	User          *schema.User
}

func (a *Transaction) TableName() string {
	return config.C.FormatTableName("transaction")
}

type Transactions []*Transaction

type QueryTransactionParam struct {
	util.PaginationParam
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	Coin        string `json:"coin"`
	Amount      int64  `json:"amount,string,omitempty"`
	TxHash      string `json:"tx_hash"`
	UserId      string `json:"user_id"`
	Status      int8   `json:"status"`
}

type TransactionQueryResult struct {
	Data       Transactions
	PageResult *util.PaginationResult
}

type TransactionQueryOptions struct {
	util.QueryOptions
}
