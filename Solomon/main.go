package main

import (
	"fastHTTP/Solomon/solomon"
	"log"
	"net/http"
)

func main() {
	slm := solomon.New()
	//非group 上加上 mw
	//slm.Use(testMidWare)
	slm.Get("/write", write)
	slm.Get("/json", json)
	slm.Get("/html", html)
	slm.Get("/str", str)
	//上面是基础返回功能测试路由
	slm.Get("/hello/*/index", allStar)
	slm.Get("/hello/he/index", test)
	slm.Get("/hello/*file/index", preStar)
	slm.Get("/hello/hel*/index", rearStar)
	slm.Get("/hello/h*ing", midStar)
	slm.Get("/hello/:name/index", peram)
	slm.Get("/hello/query", testQuery)
	//上面是基础路由的返回
	g1 := slm.NewGroup("/groupOne")
	{
		g1.Get("/hello", hello)
		g2 := g1.NewGroup("/g2")
		//g2.Use(testMidWare)
		{
			g2.Get("/hello", g2hello)
		}
	}
	//上面是分组路由的返回功能的测试
	slm.Run(":2333")
}

//中间件测试函数
func testMidWare(c *solomon.Context) {
	log.Printf("[%d]:server,and path is  %v:", 200, c.Path)
}

func hello(c *solomon.Context) {
	c.JSON(http.StatusOK, solomon.H{
		"group": "ok",
	})
}
func g2hello(c *solomon.Context) {
	c.HTML(http.StatusOK, "<h1>g2 is ok </h1>")
}

//...
func test(c *solomon.Context) {
	c.JSON(http.StatusOK, solomon.H{
		"riht": "sdf",
	})
}
func testQuery(c *solomon.Context) {
	c.JSON(http.StatusOK, solomon.H{
		"query": c.Query("name"),
	})
}
func peram(c *solomon.Context) {
	c.JSON(http.StatusOK, solomon.H{
		"peram": c.Param("name"),
	})
}
func allStar(ctx *solomon.Context) {
	ctx.JSON(http.StatusOK, solomon.H{
		"allStar": "success",
	})
}
func preStar(ctx *solomon.Context) {
	ctx.JSON(http.StatusOK, solomon.H{
		"preStar": "success",
	})
}
func rearStar(ctx *solomon.Context) {
	ctx.JSON(http.StatusOK, solomon.H{
		"rearStar": "success",
	})
}
func midStar(c *solomon.Context) {
	c.JSON(http.StatusOK, solomon.H{
		"midstar": "success",
	})
}

//...
func write(ctx *solomon.Context) {
	ctx.Write(http.StatusOK, []byte("write is ok "))
}
func json(ctx *solomon.Context) {
	ctx.JSON(http.StatusOK, solomon.H{
		"json": "success",
	})
}
func html(ctx *solomon.Context) {
	ctx.HTML(http.StatusOK, "<h1>html ok </h1>")
}
func str(ctx *solomon.Context) {
	ctx.String(http.StatusOK, "string is ok ")
}
