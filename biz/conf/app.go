package conf

import (
	"github.com/xh-polaris/synapse/biz/pkg/errorx"
	"github.com/xh-polaris/synapse/biz/types/errno"
)

type App struct {
	Status int
}

// ValidApp check whether the app valid
func ValidApp(app string) error {
	if v, ok := GetConfig().App[app]; ok {
		if v.Status != 0 {
			return errorx.New(errno.InvalidApp, errorx.KV("name", app))
		}
		return nil
	}
	return errorx.New(errno.UnSupportApp, errorx.KV("name", app))
}
