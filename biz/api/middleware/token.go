package middleware

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/xh-polaris/synapse/biz/application/base/token"
	"github.com/xh-polaris/synapse/biz/conf"
	ctxcache "github.com/xh-polaris/synapse/biz/pkg/ctxcache/ctx_cache"
	"github.com/xh-polaris/synapse/biz/types/cst"
)

func ExtractTokenInfoMW() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		str := string(ctx.GetHeader("Authorization"))
		if str != "" {
			info, err := token.ParseJWT(conf.GetConfig().Token, str)
			if err != nil {
				ctx.AbortWithMsg(err.Error(), http.StatusUnauthorized)
				return
			}
			ctxcache.Store(c, cst.TokenInfo, info)
		}
		ctx.Next(c)
	}
}
