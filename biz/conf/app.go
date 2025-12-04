package conf

import (
	"github.com/xh-polaris/synapse/biz/api/model/base"
	"github.com/xh-polaris/synapse/biz/pkg/errorx"
	"github.com/xh-polaris/synapse/biz/types/errno"
)

type App struct {
	Status    int
	ResetCode string
}

// ValidApp check whether the app valid
func ValidApp(app *base.App) error {
	if app == nil {
		return errorx.New(errno.MissingParameter, errorx.KV("parameter", "app"))
	}
	name := app.GetName()
	if v, ok := GetConfig().App[name]; ok {
		if v.Status != 0 {
			return errorx.New(errno.InvalidApp, errorx.KV("name", name))
		}
		return nil
	}
	return errorx.New(errno.UnSupportApp, errorx.KV("name", name))
}

// VerifyResetCode 检查Code是否正确以重置密码
func VerifyResetCode(app *base.App, code string) (error, bool) {
	if app == nil {
		return errorx.New(errno.MissingParameter, errorx.KV("parameter", "app")), false
	}
	name := app.GetName()
	if v, ok := GetConfig().App[name]; ok {
		return nil, code == v.ResetCode
	}
	return nil, false
}
