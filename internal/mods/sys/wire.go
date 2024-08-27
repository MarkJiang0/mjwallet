package sys

import (
	"github.com/markjiang0/mjwallet/internal/mods/sys/api"
	"github.com/markjiang0/mjwallet/internal/mods/sys/biz"
	"github.com/markjiang0/mjwallet/internal/mods/sys/dal"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	wire.Struct(new(SYS), "*"),
	wire.Struct(new(dal.Logger), "*"),
	wire.Struct(new(biz.Logger), "*"),
	wire.Struct(new(api.Logger), "*"),
)
