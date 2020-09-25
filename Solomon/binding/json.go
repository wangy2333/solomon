package binding

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type jsonBinding struct {
}

func (jb jsonBinding) Name() string {
	return "jsonBind"
}

func (jb jsonBinding) Bind(ctx *fasthttp.RequestCtx, obj interface{}) error {
	body := ctx.PostBody()
	return json.Unmarshal(body, obj)
}
