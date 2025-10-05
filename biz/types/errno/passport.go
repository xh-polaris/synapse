package errno

import "github.com/xh-polaris/synapse/biz/pkg/errorx/code"

// System 200 000 000	~ 100 999 999
const (
	UnSupportAuthType = 200_000_000
	MustPassword      = 200_000_001
	ErrVerifyCode     = 200_000_002
	PhoneHasExisted   = 200_000_003
	PhoneNotExisted   = 200_000_004
	ErrRegister       = 200_000_005
	NoPassword        = 200_000_006
	ErrPassword       = 200_000_007
)

func init() {
	code.Register(
		UnSupportAuthType,
		"the auth type {type} is not supported",
		code.WithAffectStability(false),
	)
	code.Register(
		MustPassword,
		"password is required",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrVerifyCode,
		"verify is error",
		code.WithAffectStability(false),
	)
	code.Register(
		PhoneHasExisted,
		"phone {phone} has existed",
		code.WithAffectStability(false),
	)
	code.Register(
		PhoneNotExisted,
		"phone {phone} not exists",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrRegister,
		"register error, please try again",
		code.WithAffectStability(false),
	)
	code.Register(
		NoPassword,
		"no password has been set",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrPassword,
		"error password",
		code.WithAffectStability(false),
	)
}
