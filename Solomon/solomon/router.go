package solomon

import (
	"fmt"
	"net/http"
)

type Router struct {
	handles map[string]handleFn //依附与tree存在的，
	trees   map[string]*Trie
}

//创建一个router
func NewRouter() *Router {
	return &Router{
		handles: map[string]handleFn{},
		trees:   map[string]*Trie{},
	}
}

//添加一个路由
func (r *Router) addRoute(method string, path string, handle handleFn) {
	//添加路由是先加入到树当中， 然后再放在handles当中
	if _, ok := r.trees[method]; !ok {
		//不存在这个方法树，就懒加载创建一个方法树
		r.trees[method] = NewTrie()
	}
	r.trees[method].insertNode(path) // 将path掺入对应方法树中
	key := method + ":" + path
	if r.handles == nil {
		r.handles = make(map[string]handleFn)
	}
	r.handles[key] = handle //最后我们拼接方法和路径，放到handles当中存起来对应的handle
}

//查询路由节点
func (r *Router) findRouter(method string, path string) (*routeNode, error) {

	if tree, ok := r.trees[method]; ok {
		//存在这个方法树，开始准备寻找
		node, err := tree.searchNode(path)
		if err != nil {
			return nil, fmt.Errorf("Not router node ! path is :%v ,and err:%v ", path, err)
		}
		return node, nil
	} else {
		return nil, fmt.Errorf("Not method tree! method is :%v ", method)
	}
}

func (r *Router) handleServer(c *Context) {
	if c.Path == "/favicon.ico" { //会自动访问那个小图标，挡住免得恶心人
		return
	}
	node, err := r.findRouter(c.Method, c.Path)
	//fmt.Println(node)
	if err != nil || node == nil {
		fmt.Println("server err: ", err)
		return
	}
	if len(node.params) != 0 {
		c.Parameter = node.params
	}
	key := c.Method + ":" + node.path
	if fn, ok := r.handles[key]; ok {
		c.handles = append(c.handles, fn)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
	//将所有需要执行的函数都串起来，  存在c 的handles当中，然后调用next统一执行
	c.Next()
}
