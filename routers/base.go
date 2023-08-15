package routers

import (
	"github.com/astaxie/beego/context"
	"github.com/casbin/caswaf/controllers"
)

func responseError(ctx *context.Context, error string, data ...interface{}) {
	resp := controllers.Response{Status: "error", Msg: error}
	switch len(data) {
	case 2:
		resp.Data2 = data[1]
		fallthrough
	case 1:
		resp.Data = data[0]
	}

	err := ctx.Output.JSON(resp, true, false)
	if err != nil {
		panic(err)
	}
}

func denyRequest(ctx *context.Context) {
	responseError(ctx, "Unauthorized operation")
}
