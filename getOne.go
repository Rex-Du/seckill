// Author : rexdu
// Time : 2020-04-04 00:32
package main

import (
	"github.com/kataras/golog"
	"net/http"
	"sync"
)

var sum int64 = 0

var total int64 = 100000

//var req_count = 0 // 抢购的人数，用来控制放量，比如每100个人请求，才放出一个商品

var mutex sync.Mutex

func GetOneProduct() bool {
	mutex.Lock()
	defer mutex.Unlock()
	//req_count += 1
	//if req_count%100 == 0 {
	if sum < total {
		sum += 1
		golog.Println("已售出：", sum)
		return true
		//}
	}
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
