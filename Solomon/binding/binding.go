package binding

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

type Binder interface {
	Name() string                                             //当前准备的方法名称
	Bind(request *fasthttp.RequestCtx, obj interface{}) error //所有类型的入口函数
}

//常见的几种基本类型
const (
	JSON = "application/json"
	//MIMEHTML              = "text/html"
	//MIMEXML               = "application/xml"
	//MIMEXML2              = "text/xml"
	//MIMEPlain             = "text/plain"
	//MIMEMultipartPOSTForm = "multipart/form-data"
)

//不同的类型对应的不同对象， 对象当中包含各自处理函数的方法 都满足这个binder接口
var (
	JsonType = jsonBinding{}
	Form     = formBinding{}
	//...
)

//创建的bind 需要当前请求的方法， 和content-type
func NewBind(method string, contentType string) Binder {
	if method == http.MethodGet {
		return Form
	}
	switch contentType {
	case JSON:
		return JsonType
	default:
		return Form
	}
}
