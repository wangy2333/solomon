package solomon

import (
	"strings"

	"github.com/valyala/fasthttp"
)

//solomon的引擎
type Engine struct {
	router *Router
	*solomonGroup
	groups []*solomonGroup
}

//分组
type solomonGroup struct {
	prefix string
	//groupName string
	engine  *Engine
	parent  *solomonGroup
	midWare []handleFn
}
type handleFn func(ctx *Context)

//创建一个solomon的实例
func New() *Engine {
	engine := &Engine{
		router: NewRouter(),
	}
	engine.solomonGroup = &solomonGroup{engine: engine}
	engine.groups = []*solomonGroup{engine.solomonGroup}
	return engine
}

//创建一个solomonGroup的实例
func (group *solomonGroup) NewGroup(groupName string) *solomonGroup {
	engine := group.engine
	newGroup := &solomonGroup{
		prefix: group.prefix + groupName, //可以组内分组construct
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

//midWare 的方法
func (group *solomonGroup) Use(fn ...handleFn) {
	group.midWare = append(group.midWare, fn...)
}

//基本四大方法
//get func
func (group *solomonGroup) Get(path string, handle handleFn) {
	PATH := group.prefix + path
	group.engine.router.addRoute("GET", PATH, handle)
}

//post func
func (group *solomonGroup) Post(path string, handle handleFn) {
	PATH := group.prefix + path
	group.engine.router.addRoute("POST", PATH, handle)
}

//post func
func (group *solomonGroup) Delete(path string, handle handleFn) {
	PATH := group.prefix + path
	group.engine.router.addRoute("DELETE", PATH, handle)
}

//post func
func (group *solomonGroup) Put(path string, handle handleFn) {
	PATH := group.prefix + path
	group.engine.router.addRoute("PUT", PATH, handle)
}

//所有方法的入口函数
//得到了方法和路径之后， 先去handleTrees当中寻找，
//然后找到了就得到path而不是用ctx中的path
//因为ctx的path当中可能带有通配符和参数，直接查询是没有意义的
//而路由树当中的记录和handles才是一一对应的。 所以我们拿着ctx.path去找path
func (e *Engine) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	c := newContext(ctx)
	var midWare []handleFn
	//遍历组，将组的前缀和请求的路径进行比对， 将符合当前请求的组的中间件放入其中，
	//不是当前的中间件不执行
	for _, group := range e.groups {
		if strings.HasPrefix(c.Path, group.prefix) {
			midWare = append(midWare, group.midWare...)
		}
	}
	//准备好中间件之后就可以传给c进行执行， 而不是在这里执行，
	//因为c由函数可以去遍历中间件进行handles
	c.handles = midWare
	e.router.handleServer(c)
}

//能够监听
func (e *Engine) Run(address string) {
	fasthttp.ListenAndServe(address, e.HandleFastHTTP)
}
