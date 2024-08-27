package schema

import (
	"github.com/markjiang0/mjwallet/internal/config"
	"github.com/markjiang0/mjwallet/internal/mods/rbac/schema"
	"github.com/markjiang0/mjwallet/pkg/util"
	"time"
)

type Account struct {
	ID            string    `json:"id" gorm:"size:20;primarykey;"`
	UserId        string    `json:"user_id" gorm:"size:20"`
	Address       string    `json:"address" gorm:"size:128"`
	AddressBase58 string    `json:"address_base_58" gorm:"size:64"`
	PrivateKey    string    `json:"private_key" gorm:"size:2048"`
	PublicKey     string    `json:"public_key" gorm:"size:2048"`
	Balance       int64     `json:"balance"`
	UsdtBalance   int64     `json:"usdt_balance"`
	CreatedAt     time.Time `json:"created_at" gorm:"index;"` // Create time
	UpdatedAt     time.Time `json:"updated_at" gorm:"index;"` // Update time
	Default       int8      `json:"default"`
	Name          string    `json:"name" gorm:"size:50"`
	User          *schema.User
}

type AccountDto struct {
	ID            string    `json:"id"`
	UserId        string    `json:"user_id"`
	Username      string    `json:"username"`
	AddressBase58 string    `json:"address_base_58"`
	Balance       int64     `json:"balance"`
	UsdtBalance   int64     `json:"usdt_balance"`
	Default       bool      `json:"default"`
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"` // Create time
}

type TransferParam struct {
	ToAddress string `json:"to_address"`
	Amount    int64  `json:"amount"`
}

type QueryAccountParam struct {
	util.PaginationParam
	LikeUserName string `form:"user_name"`
	UserId       string `form:"user_id"`
	Default      string `form:"default"`
}

type AccountQueryResult struct {
	Data       AccountDtos
	PageResult *util.PaginationResult
}

type AccountQueryOptions struct {
	util.QueryOptions
}

type Accounts []*Account
type AccountDtos []*AccountDto

func (a *Account) TableName() string {
	return config.C.FormatTableName("account")
}
