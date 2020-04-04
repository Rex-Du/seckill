// Author : rexdu
// Time : 2020-04-02 22:33
package common

import (
	"net/http"
	"strings"
)

type FilterHandle func(rw http.ResponseWriter, req *http.Request) error

type Filter struct {
	filterMap map[string]FilterHandle
}

func NewFilter() *Filter {
	return &Filter{filterMap: make(map[string]FilterHandle)}
}

func (f *Filter) RegisterFilterUri(uri string, handler FilterHandle) {
	f.filterMap[uri] = handler
}

func (f *Filter) GetFilterHandler(uri string) FilterHandle {
	return f.filterMap[uri]
}

type WebHandle func(rw http.ResponseWriter, req *http.Request)

// 执行拦截器，返回业务类型
func (f *Filter) Handle(webHandle WebHandle) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		for path, handle := range f.filterMap {
			//if path == r.RequestURI {
			if strings.Contains(r.RequestURI, path) {
				// 执行拦截业务逻辑
				err := handle(rw, r)
				if err != nil {
					rw.Write([]byte(err.Error()))
					return
				}
				break
			}
		}
		// 执行正常注册的函数
		webHandle(rw, r)
	}
}
