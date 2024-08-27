package wallet

import (
	"github.com/google/wire"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/api"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/biz"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/dal"
	"github.com/markjiang0/mjwallet/internal/mods/wallet/task"
)

var Set = wire.NewSet(
	wire.Struct(new(WALLET), "*"),
	wire.Struct(new(dal.Account), "*"),
	wire.Struct(new(biz.Account), "*"),
	wire.Struct(new(api.Account), "*"),
	wire.Struct(new(task.Account), "*"),
	wire.Struct(new(dal.Transaction), "*"),
	wire.Struct(new(biz.Transaction), "*"),
	wire.Struct(new(api.Transaction), "*"),
	wire.Struct(new(task.Transaction), "*"),
	wire.Struct(new(biz.TelegramBot), "*"),
	wire.Struct(new(dal.TelegramBot), "*"),
	wire.Struct(new(api.TelegramBot), "*"),
	wire.Struct(new(task.TelegramBot), "*"),
)
