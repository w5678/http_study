package main

import (
	"fmt"
	"log"
	"net/http"
)

/*
构造了空结构体来实现了 ServeHTTP这个接口
这个接口:
第一个参数是 w response，利用 ResponseWrite来构造的response对象
第二个参数是 req
*/

type Engine struct{}

func (eng *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path=%q\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 Not Found : %s\n", req.URL)

	}
}

func main() {	
	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":9999", engine))
}
