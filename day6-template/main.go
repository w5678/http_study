package main

import (
	"fmt"
	"gee"
	template "html/template"
	"net/http"
	"time"
)

/*
middleware 指的是非业务的技术组件
框架需要一个插口，允许用户自己定义功能，嵌入到框架中去
关键点是：
插入点在哪
中间件的输入是什么
*/

//func onlyForV2() gee.HandlerFunc {
//	return func(c *gee.Context) {
//		t := time.Now()
//		//c.Fail(500,"Internal Server Error")
//		log.Printf("[%d] %s in %v for group v2-1", c.StatusCode, c.Req.RequestURI, time.Since(t))
//	}
//}
//
//func onlyForV21() gee.HandlerFunc {
//	return func(c *gee.Context) {
//		t := time.Now()
//		//c.Fail(500,"Internal Server Error")
//		log.Printf("[%d] %s in %v for group v2-2", c.StatusCode, c.Req.RequestURI, time.Since(t))
//	}
//}
//
//
//func main() {
//	r := gee.New()
//	r.Use(gee.Logger())
//	r.GET("/index", func(c *gee.Context) {
//		c.HTML(http.StatusOK, "<h1>Hello index</h1>")
//	})
//	v1 := r.Group("/v1")
//
//	v1.GET("/", func(c *gee.Context) {
//		c.HTML(http.StatusOK, "<h1>Hello gee</h1>")
//	})
//	v1.GET("/hello", func(c *gee.Context) {
//		c.String(http.StatusOK, "hello %s,you are at %s\n", c.Query("name"), c.Path)
//	})
//	v2 := r.Group("/v2")
//	v2.Use(onlyForV2(),onlyForV21())
//	{
//		v2.GET("/hello/:name", func(c *gee.Context) {
//			c.String(http.StatusOK, "hello %s,you are at %s\n", c.Param("name"), c.Path)
//		})
//		v2.POST("/login", func(c *gee.Context) {
//			c.JSON(http.StatusOK, gee.H{
//				"username": c.PostForm("username"),
//				"password": c.PostForm("password"),
//			})
//		})
//	}
//	r.Run(":9999")
//
//}

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := gee.New()
	r.Use(gee.Logger())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *gee.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gee.H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *gee.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
			"title": "gee",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	r.Run(":9999")
}
