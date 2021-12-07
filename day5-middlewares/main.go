package main

/*
middleware 指的是非业务的技术组件
框架需要一个插口，允许用户自己定义功能，嵌入到框架中去
关键点是：
插入点在哪
中间件的输入是什么
*/

//导入gee自定义模块和 系统模块http
import (
	"gee"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		t := time.Now()
		//c.Fail(500,"Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := gee.New()
	r.Use(gee.Logger())
	r.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello index</h1>")
	})
	v1 := r.Group("/v1")

	v1.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello gee</h1>")
	})
	v1.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s,you are at %s\n", c.Query("name"), c.Path)
	})
	v2 := r.Group("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s,you are at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}
	r.Run(":9999")

}
