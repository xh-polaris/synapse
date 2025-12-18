package errno

import "github.com/xh-polaris/synapse/biz/pkg/errorx/code"

// System 200 000 000	~ 100 999 999
const (
	UnSupportAuthType   = 200_000_000
	MustPassword        = 200_000_001
	ErrVerifyCode       = 200_000_002
	PhoneHasExisted     = 200_000_003
	PhoneNotExisted     = 200_000_004
	ErrRegister         = 200_000_005
	NoPassword          = 200_000_006
	ErrPassword         = 200_000_007
	ErrResetPassword    = 200_000_008
	UnSupportThirdParty = 200_000_009
	ErrThirdPartyLogin  = 200_000_010
	TooOftenLoginError  = 200_000_011
	CodeHasExisted      = 200_000_012
	CodeNotExisted      = 200_000_013
	MissingParameter    = 200_000_014
	InvalidToken        = 200_000_015
	UserNotExisted      = 200_000_016
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
		"phone not exists",
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
	code.Register(
		ErrResetPassword,
		"set password failed",
		code.WithAffectStability(false),
	)
	code.Register(UnSupportThirdParty,
		"unsupported third party",
		code.WithAffectStability(false),
	)
	code.Register(ErrThirdPartyLogin,
		"third party login failed",
		code.WithAffectStability(false),
	)
	code.Register(TooOftenLoginError,
		"登录失败次数过多, 请 {period} 分钟后再试",
		code.WithAffectStability(false),
	)
	code.Register(CodeHasExisted,
		"student id {code} has existed",
		code.WithAffectStability(false),
	)
	code.Register(CodeNotExisted,
		"student id {code} not exists",
		code.WithAffectStability(false),
	)
	code.Register(MissingParameter,
		"missing parameter {parameter}",
		code.WithAffectStability(false),
	)
	code.Register(InvalidToken,
		"invalid token",
		code.WithAffectStability(false),
	)
}
