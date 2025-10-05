package internal

import "github.com/xh-polaris/synapse/biz/api/model/base"

func Success() *base.Response {
	return &base.Response{Code: 0, Msg: ""}
}
