package common

import (
	"net/http"
	"strings"
)

// 声明一个新的数据类型（函数类型）
type FilterHandler func(respW http.ResponseWriter, req *http.Request) error

type WebHandler func(respW http.ResponseWriter, req *http.Request)

// 拦截器结构体
type Filter struct {
	// 用来存储需要拦截的 URI
	filterUriMap map[string]FilterHandler
}

// 拦截器初始化
func NewFilter() *Filter {
	return &Filter{
		filterUriMap: make(map[string]FilterHandler),
	}
}

// 注册拦截器
func (f *Filter) RegisterUriFilter(uri string, filterHandler FilterHandler) {
	f.filterUriMap[uri] = filterHandler
}

// 获取拦截器对应的 Handler
func (f *Filter) GetFilter(uri string) FilterHandler {
	return f.filterUriMap[uri]
}

// 执行拦截器，返回函数类型
func (f *Filter) Handler(webHandler WebHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for path, handle := range f.filterUriMap {
			if strings.Contains(r.RequestURI, path) {
				err := handle(w, r)
				if err != nil {
					w.Write([]byte(err.Error()))
					return
				}
				break
			}
		}
		// 执行正常注册函数
		webHandler(w, r)
	}

}
