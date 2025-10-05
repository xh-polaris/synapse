package synapse

import (
	"context"
	"reflect"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server/render"
	"github.com/xh-polaris/synapse/biz/api/internal/httputil"
)

func init() {
	SetJSONMarshal()
}

func invalidParamRequestResponse(c *app.RequestContext, errMsg string) {
	httputil.BadRequest(c, errMsg)
}

func internalServerErrorResponse(ctx context.Context, c *app.RequestContext, err error) {
	httputil.InternalError(ctx, c, err)
}

func SetJSONMarshal() {
	render.ResetJSONMarshal(func(v interface{}) (data []byte, err error) {
		if data, err = sonic.ConfigFastest.Marshal(v); err != nil {
			return nil, err
		}
		if !hasRespField(v) {
			return data, nil
		}

		var root ast.Node
		var node *ast.Node
		if root, err = sonic.Get(data); err != nil {
			return nil, err
		}
		if node = root.Get("resp"); node.Check() != nil {
			return data, nil
		}

		if _, err = root.Set("code", *node.Get("code")); err != nil {
			return nil, err
		}
		if _, err = root.Set("msg", *node.Get("msg")); err != nil {
			return nil, err
		}
		if extra := node.Get("extra"); extra != nil {
			if _, err = root.Set("extra", *extra); err != nil {
				return nil, err
			}
		}
		if _, err = root.Unset("resp"); err != nil {
			return nil, err
		}
		return root.MarshalJSON()
	})
}

func hasRespField(v interface{}) bool {
	rt := reflect.TypeOf(v)
	for rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Struct || rt.NumField() == 0 {
		return false
	}
	firstField := rt.Field(0)
	return firstField.Name == "Resp" &&
		firstField.Type.String() == "*base.Response"
}
