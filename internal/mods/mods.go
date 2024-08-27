package mods

import (
	"context"
	"github.com/markjiang0/mjwallet/internal/mods/wallet"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/markjiang0/mjwallet/internal/mods/rbac"
	"github.com/markjiang0/mjwallet/internal/mods/sys"
)

const (
	apiPrefix = "/api/"
)

// Collection of wire providers
var Set = wire.NewSet(
	wire.Struct(new(Mods), "*"),
	rbac.Set,
	sys.Set,
	wallet.Set,
)

type Mods struct {
	RBAC   *rbac.RBAC
	SYS    *sys.SYS
	WALLET *wallet.WALLET
}

func (a *Mods) Init(ctx context.Context) error {
	if err := a.RBAC.Init(ctx); err != nil {
		return err
	}
	if err := a.SYS.Init(ctx); err != nil {
		return err
	}
	if err := a.WALLET.Init(ctx); err != nil {
		return err
	}

	return nil
}

func (a *Mods) RouterPrefixes() []string {
	return []string{
		apiPrefix,
	}
}

func (a *Mods) RegisterRouters(ctx context.Context, e *gin.Engine) error {
	gAPI := e.Group(apiPrefix)
	v1 := gAPI.Group("v1")

	if err := a.RBAC.RegisterV1Routers(ctx, v1); err != nil {
		return err
	}
	if err := a.SYS.RegisterV1Routers(ctx, v1); err != nil {
		return err
	}
	if err := a.WALLET.RegisterV1Routers(ctx, v1); err != nil {
		return err
	}

	return nil
}

func (a *Mods) Release(ctx context.Context) error {
	if err := a.RBAC.Release(ctx); err != nil {
		return err
	}
	if err := a.SYS.Release(ctx); err != nil {
		return err
	}
	if err := a.WALLET.Release(ctx); err != nil {
		return err
	}
	return nil
}
