package solomon

import (
	"encoding/json"
	"fastHTTP/Solomon/binding"
	"fmt"
	"unsafe"

	"github.com/valyala/fasthttp"
)

type H map[string]interface{}
type Context struct {
	f           *fasthttp.RequestCtx
	StatusCode  int
	Path        string
	Method      string
	ContentType string
	Parameter   map[string]string
	//midWare
	handles []handleFn
	index   int
}

func newContext(f *fasthttp.RequestCtx) *Context {
	return &Context{
		f:           f,
		Path:        byteToString(f.Path()),
		Method:      byteToString(f.Method()),
		ContentType: byteToString(f.Request.Header.ContentType()),
		Parameter:   make(map[string]string),
		handles:     make([]handleFn, 0),
		index:       -1,
	}
}

//midWare function
//这个函数本省就是为了遍历c当中的函数组，
//之所以引入index就是为了解决一个中间件执行一半可以
//调用该中间件的next然后执行中间件的中间件，执行完毕再跳回来继续遍历当层
func (c *Context) Next() {
	c.index++
	for ; c.index < len(c.handles); c.index++ {
		c.handles[c.index](c)
	}
}

//[]byte to string
func byteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//set statusCode
func (c *Context) status(code int) {
	c.StatusCode = code
	c.f.Response.Header.SetStatusCode(code)
}

//set head key and value
func (c *Context) setHead(key, value string) {
	c.f.Response.Header.Set(key, value)
}

//error show to client
func (c *Context) error(status int, err string) {
	c.status(status)
	c.f.Error(err, status)
}

//write to byte
func (c *Context) Write(statusCode int, b []byte) {
	c.status(statusCode)
	c.f.Write(b)
}

//write to json
func (c *Context) JSON(statusCode int, obj interface{}) {
	c.status(statusCode)
	c.setHead("Content-Type", "application/json")
	encoder := json.NewEncoder(c.f) //将fastHttp的write创建一个encoder然后给obj
	err := encoder.Encode(obj)
	if err != nil {
		c.error(500, err.Error())
	}
}

//write to string
func (c *Context) String(statusCode int, format string, content ...interface{}) {
	c.status(statusCode)
	c.setHead("Content-Type", "text/plain")
	c.f.Write([]byte(fmt.Sprintf(format, content...)))
}
func (c *Context) HTML(status int, html string) {
	c.status(status)
	c.setHead("Content-Type", "text/html")
	c.f.Write([]byte(html))
}

//将url中的参数拿出来
func (c *Context) Query(key string) string {
	return byteToString(c.f.QueryArgs().Peek(key))
}

//查询uel中带的key是否存在value 如果存在就返回value
func (c *Context) Param(key string) string {
	value, _ := c.Parameter[key]
	return value
}

//context 的shouldBind函数
//根据传入的参数和contentType返回不同的实例调用接口中的方法bind
func (c *Context) ShouldBind(obj interface{}) error {
	b := binding.NewBind(c.Method, c.ContentType)
	return c.shouldBindWith(obj, b)
}
func (c *Context) shouldBindWith(obj interface{}, b binding.Binder) error {
	return b.Bind(c.f, obj)
}
