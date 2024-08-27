package bootstrap

import (
	"context"
	"github.com/markjiang0/mjwallet/internal/wirex"
)

func startCronJob(ctx context.Context, injector *wirex.Injector) error {
	err := injector.M.WALLET.RegisterAndStartCronJob(ctx)
	return err
}
