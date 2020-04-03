// Author : rexdu
// Time : 2020-04-04 00:32
package main

import (
	"github.com/kataras/golog"
	"net/http"
	"sync"
)

var sum int64 = 0

var total int64 = 10000

var mutex sync.Mutex

func GetOneProduct() bool {
	mutex.Lock()
	defer mutex.Unlock()
	if sum < total {
		sum += 1
		return true
	}
	golog.Info("已售罄！")
	return false
}

func GetProduct(rw http.ResponseWriter, req *http.Request) {
	if GetOneProduct() {
		rw.Write([]byte("true"))
		return
	}
	rw.Write([]byte("false"))
	return
}

func main() {
	http.HandleFunc("/getOne", GetProduct)
	http.ListenAndServe(":8080", nil)
}
