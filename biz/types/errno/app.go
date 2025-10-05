package errno

import "github.com/xh-polaris/synapse/biz/pkg/errorx/code"

// App 101 000 000	~ 101 999 999
const (
	InvalidApp   = 101_000_000
	UnSupportApp = 101_000_001
)

func init() {
	code.Register(
		InvalidApp,
		"the app {name} is invalid",
		code.WithAffectStability(false),
	)
	code.Register(
		UnSupportApp,
		"the app {name} is not supported",
		code.WithAffectStability(false),
	)
}
