package errno

import "github.com/xh-polaris/synapse/biz/pkg/errorx/code"

// System 100 000 000	~ 100 999 999
const (
	ErrInvalidAuthType = 100_000_000
	ErrSendPhoneVerify = 100_000_001
)

func init() {
	code.Register(
		ErrInvalidAuthType,
		"the auth type is invalid",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrSendPhoneVerify,
		"send phone verify code failed",
		code.WithAffectStability(false),
	)
}
