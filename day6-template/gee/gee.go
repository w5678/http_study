package gee

/*
对原生对http模块进行了封装
*/
import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

//定义一个函数类型，用于回调
type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc //support middleware
	parent      *RouterGroup  //support nesting
	engine      *Engine       //all groups share a Engine instance
}

//创建引擎类/结构体
type Engine struct {
	*RouterGroup
	router        *router            //router也是个自定义的router，私有
	groups        []*RouterGroup     //stores all groups
	htmlTemplates *template.Template //for html render
	funcMap       template.FuncMap   //for html render
}

//new返回一engine实例，router也被创建出来
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

//engine的方法
func (group *RouterGroup) addRouter(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s -%s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// engine的方法，在get的基础上加入定制化的handle函数
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRouter("GET", pattern, handler) //组合能力，将engine的方法给继承过来了
}
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRouter("POST", pattern, handler)
}

//启动http serverHandlerFunc
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

//将middlewares加入到group中
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

//create static handler

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		//check not exist or no access permission
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root)) // create a static file handler func
	urlPattern := path.Join(relativePath, "/*filepath")                //concat the url pattern
	group.GET(urlPattern, handler)                                     //register
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

//启动server时候，通过前缀来判断需要启用哪些个middleware，得到中间件列表后赋值给handlers
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}
