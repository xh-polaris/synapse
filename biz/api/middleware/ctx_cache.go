package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	ctxcache "github.com/xh-polaris/synapse/biz/pkg/ctxcache/ctx_cache"
)

func ContextCacheMW() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		c = ctxcache.Init(c)
		ctx.Next(c)
	}
}
