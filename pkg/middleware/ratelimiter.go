package middleware

import (
	"context"
	"time"

	"github.com/markjiang0/mjwallet/pkg/errors"
	"github.com/markjiang0/mjwallet/pkg/logging"
	"github.com/markjiang0/mjwallet/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type RateLimiterConfig struct {
	Enable              bool
	AllowedPathPrefixes []string
	SkippedPathPrefixes []string
	Period              int
	MaxRequestsPerIP    int
	MaxRequestsPerUser  int
	StoreType           string // memory/redis
	MemoryStoreConfig   RateLimiterMemoryConfig
	RedisStoreConfig    RateLimiterRedisConfig
}

func RateLimiterWithConfig(config RateLimiterConfig) gin.HandlerFunc {
	if !config.Enable {
		return Empty()
	}

	var store RateLimiterStorer
	switch config.StoreType {
	case "redis":
		store = NewRateLimiterRedisStore(config.RedisStoreConfig)
	default:
		store = NewRateLimiterMemoryStore(config.MemoryStoreConfig)
	}

	return func(c *gin.Context) {
		if !AllowedPathPrefixes(c, config.AllowedPathPrefixes...) ||
			SkippedPathPrefixes(c, config.SkippedPathPrefixes...) {
			c.Next()
			return
		}

		var (
			allowed bool
			err     error
		)

		ctx := c.Request.Context()
		if userID := util.FromUserID(ctx); userID != "" {
			allowed, err = store.Allow(ctx, userID, time.Second*time.Duration(config.Period), config.MaxRequestsPerUser)
		} else {
			allowed, err = store.Allow(ctx, c.ClientIP(), time.Second*time.Duration(config.Period), config.MaxRequestsPerIP)
		}

		if err != nil {
			logging.Context(ctx).Error("Rate limiter middleware error", zap.Error(err))
			util.ResError(c, errors.InternalServerError("", "Internal server error, please try again later."))
		} else if allowed {
			c.Next()
		} else {
			util.ResError(c, errors.TooManyRequests("", "Too many requests, please try again later."))
		}
	}
}

type RateLimiterStorer interface {
	Allow(ctx context.Context, identifier string, period time.Duration, maxRequests int) (bool, error)
}

func NewRateLimiterMemoryStore(config RateLimiterMemoryConfig) RateLimiterStorer {
	return &RateLimiterMemoryStore{
		cache: cache.New(config.Expiration, config.CleanupInterval),
	}
}

type RateLimiterMemoryConfig struct {
	Expiration      time.Duration
	CleanupInterval time.Duration
}

type RateLimiterMemoryStore struct {
	cache *cache.Cache
}

func (s *RateLimiterMemoryStore) Allow(ctx context.Context, identifier string, period time.Duration, maxRequests int) (bool, error) {
	if period.Seconds() <= 0 || maxRequests <= 0 {
		return true, nil
	}

	if limiter, exists := s.cache.Get(identifier); exists {
		isAllow := limiter.(*rate.Limiter).Allow()
		s.cache.SetDefault(identifier, limiter)
		return isAllow, nil
	}

	limiter := rate.NewLimiter(rate.Every(period), maxRequests)
	limiter.Allow()
	s.cache.SetDefault(identifier, limiter)

	return true, nil
}

type RateLimiterRedisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}

func NewRateLimiterRedisStore(config RateLimiterRedisConfig) RateLimiterStorer {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Username: config.Username,
		Password: config.Password,
		DB:       config.DB,
	})

	return &RateLimiterRedisStore{
		limiter: redis_rate.NewLimiter(rdb),
	}
}

type RateLimiterRedisStore struct {
	limiter *redis_rate.Limiter
}

func (s *RateLimiterRedisStore) Allow(ctx context.Context, identifier string, period time.Duration, maxRequests int) (bool, error) {
	if period.Seconds() <= 0 || maxRequests <= 0 {
		return true, nil
	}

	result, err := s.limiter.Allow(ctx, identifier, redis_rate.PerSecond(maxRequests/int(period.Seconds())))
	if err != nil {
		return false, err
	}
	return result.Allowed > 0, nil
}
