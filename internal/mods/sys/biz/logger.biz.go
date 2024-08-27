package biz

import (
	"context"

	"github.com/markjiang0/mjwallet/internal/mods/sys/dal"
	"github.com/markjiang0/mjwallet/internal/mods/sys/schema"
	"github.com/markjiang0/mjwallet/pkg/util"
)

// Logger management
type Logger struct {
	LoggerDAL *dal.Logger
}

// Query loggers from the data access object based on the provided parameters and options.
func (a *Logger) Query(ctx context.Context, params schema.LoggerQueryParam) (*schema.LoggerQueryResult, error) {
	params.Pagination = true

	result, err := a.LoggerDAL.Query(ctx, params, schema.LoggerQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: []util.OrderByParam{
				{Field: "created_at", Direction: util.DESC},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
