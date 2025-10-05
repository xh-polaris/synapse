package sms

import (
	"context"
	"fmt"
	"time"

	"github.com/xh-polaris/synapse/biz/infra/contract/cache"
)

const (
	Passport = "passport"
)

// Provider 提供SMS相关服务, 支持附有效期(param.Expire)的验证码(param.Code)存储与校验
type Provider interface {
	// Send 发送验证码
	Send(ctx context.Context, app, cause, phone string, param *SMSParam) error
	// Check 校验验证码
	Check(ctx context.Context, app, cause, phone, code string) (bool, error)
}

type SMSParam struct {
	Code   string
	Expire time.Duration
}

// Cache 提供SMS相关验证码存储与获取能力
type Cache struct {
	cache cache.Cmdable
}

func NewSMSCache(_ context.Context, cli cache.Cmdable) *Cache {
	return &Cache{cache: cli}
}

func (c *Cache) buildKey(app, cause, phone string) string {
	return fmt.Sprintf("%s:%s:%s", app, cause, phone)
}

// Store 存储验证码
func (c *Cache) Store(ctx context.Context, app, cause, phone, code string, expire time.Duration) error {
	result := c.cache.Set(ctx, c.buildKey(app, cause, phone), code, expire)
	return result.Err()
}

// Load 获取验证码
func (c *Cache) Load(ctx context.Context, app, cause, phone string) (string, error) {
	result := c.cache.Get(ctx, c.buildKey(app, cause, phone))
	return result.Val(), result.Err()
}
