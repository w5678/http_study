package gee

/*
对原生对http模块进行了封装
*/
import (
	"net/http"
)

//定义一个函数类型，用于回调
type HandlerFunc func(*Context)

//创建引擎类/结构体
type Engine struct {
	router *router //router也是个自定义的router，私有
}

//new返回一engine实例，router也被创建出来
func New() *Engine {
	return &Engine{router: newRouter()}
}

//engine的方法
func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// engine的方法，在get的基础上加入定制化的handle函数
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRouter("GET", pattern, handler)
}
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRouter("POST", pattern, handler)
}

//启动http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
