package errno

import "github.com/xh-polaris/synapse/biz/pkg/errorx/code"

// System 100 000 000	~ 100 999 999
const (
	ErrInvalidAuthType            = 100_000_000
	ErrSendPhoneVerify            = 100_000_001
	ErrSendUpperLimit             = 100_000_002
	ErrSignTicket                 = 100_000_003
	ErrExchangeTicket             = 100_000_004
	ErrExchangeTicketByInvalidKey = 100_000_005
)

func init() {
	code.Register(
		ErrInvalidAuthType,
		"the auth type {type} is invalid",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrSendPhoneVerify,
		"send phone verify code failed",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrSendUpperLimit,
		"send too often, please wait for {period} minutes")
	code.Register(
		ErrSignTicket,
		"sign ticket failed",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrExchangeTicket,
		"exchange ticket failed",
		code.WithAffectStability(false),
	)
	code.Register(
		ErrExchangeTicketByInvalidKey,
		"exchange ticket failed caused by invalid key",
		code.WithAffectStability(false),
	)
}
