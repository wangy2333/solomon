package binding

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type formBinding struct {
}

func (fb formBinding) Name() string {
	return "formBind"
}

func (fb formBinding) Bind(ctx *fasthttp.RequestCtx, obj interface{}) error {
	//先将map中的值拿出来，构建一个新的map 然后在转换到obj当中，
	form, _ := ctx.MultipartForm()
	value := form.Value
	tempMap := make(map[string]interface{})
	for k, v := range value {
		tempMap[k] = v[0]
	}
	marshal, _ := json.Marshal(tempMap)
	return json.Unmarshal(marshal, obj)
}
