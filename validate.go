// Author : rexdu
// Time : 2020-04-02 22:43
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/golog"
	"io/ioutil"
	"net/http"
	"net/url"
	"seckill/common"
	"seckill/datamodels"
	"seckill/encrypt"
	"seckill/rabbitmq"
	"strconv"
	"sync"
)

// 统一验证拦截器，每个接口都需要提前验证
func Auth(w http.ResponseWriter, r *http.Request) error {
	err := CheckUserInfo(r)
	if err != nil {
		return err
	}
	return nil
}

func CheckUserInfo(r *http.Request) error {
	// 获取uid cookie
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		return errors.New("uid获取失败")
	}
	// 获取用户加密串
	signCookie, err := r.Cookie("sign")
	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		return errors.New("加密串已被篡改！")
	}
	fmt.Println("结果比对")
	fmt.Println("用户ID：" + uidCookie.Value)
	fmt.Println("解密后用户ID：", string(signByte))
	if checkInfo(uidCookie.Value, string(signByte)) {
		return nil
	}
	return errors.New("身份校验失败！")
}

func checkInfo(checkStr string, signStr string) bool {
	if checkStr == signStr {
		return true
	}
	return false

}

func CheckRight(rw http.ResponseWriter, r *http.Request) {
	if !accessControl.GetDistributedRight(r) {
		rw.Write([]byte("false"))
		return
	}
	rw.Write([]byte("true"))
	return
}

func Check(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("执行check！")
	golog.Debug("执行check")
	// 获取url中的参数
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["productID"]) <= 0 {
		rw.Write([]byte("false"))
		return
	}
	productIDstr := queryForm["productID"][0]
	golog.Debug("获取到productID", productIDstr)
	// 获取cookie
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		rw.Write([]byte("false"))
		return
	}
	// 分布式权限验证:去访问对应的主机，
	//1.分布式权限验证
	right := accessControl.GetDistributedRight(r)
	if right == false {
		rw.Write([]byte("false"))
		return
	}
	// 获取数量控制权限，防止秒杀超卖
	hostUrl := "http://" + GetOneIp + ":" + GetOnePort + "/getOne"
	resp, body, err := GetCurl(hostUrl, r)
	if err != nil {
		rw.Write([]byte("false"))
		return
	}
	if resp.StatusCode == 200 {
		if string(body) == "true" {
			// 整合下单
			productID, err := strconv.ParseInt(productIDstr, 10, 64)
			if err != nil {
				golog.Error(productID, "失败")
				rw.Write([]byte("false"))
				return
			}
			userID, err := strconv.ParseInt(uidCookie.Value, 10, 64)
			if err != nil {
				golog.Error(userID, "失败")
				rw.Write([]byte("false"))
				return
			}
			message := datamodels.Message{productID, userID}
			msgByte, err := json.Marshal(message)
			if err != nil {
				golog.Error(message, "失败")
				rw.Write([]byte("false"))
				return
			}
			err = rabbitmqValidate.PublishSimple(string(msgByte))
			if err != nil {
				golog.Error(err, "失败")
				rw.Write([]byte("false"))
				return
			}
			rw.Write([]byte("true"))
		}
	}
	rw.Write([]byte("false"))
	return
}

// 设置集群地址
//var hostArray = []string{"192.168.124.135", "192.168.124.136"}
var hostArray = []string{"127.0.0.1"}
var localHost = "127.0.0.1"
var port = "8083"

// 一个全局一致性实例
var hashConsistent *common.Consistent
var rabbitmqValidate *rabbitmq.RabbitMQ

// 数量控制服务器内网ip
var GetOneIp = "111.229.61.201"
var GetOnePort = "8080"

// 用来存放控制信息
type AccessControl struct {
	// 存放用户想要存放的信息
	sourcesArray map[int]interface{}
	sync.RWMutex
}

var accessControl = &AccessControl{
	sourcesArray: make(map[int]interface{}),
}

// 获取指定的数据
func (m *AccessControl) GetNewRecord(uid int) interface{} {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	data := m.sourcesArray[uid]
	return data
}

//设置记录
func (m *AccessControl) SetNewRecord(uid int) {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	m.sourcesArray[uid] = "hello imooc"
}

func (m *AccessControl) GetDistributedRight(req *http.Request) bool {
	uidCookie, err := req.Cookie("uid")
	if err != nil {
		return false
	}

	// 采用一致性hash算法，确定改用户应该访问的机器
	hostRequest, err := hashConsistent.Get(uidCookie.Value)
	if err != nil {
		golog.Error("获取节点时出错", err)

		return false
	}

	uid, err := strconv.Atoi(uidCookie.Value)
	if err != nil {
		golog.Error("uid格式错误", err)
		return false
	}
	// 判断是否是本机
	if hostRequest == localHost {
		// 执行梧桐数据读取和校验
		return m.GetDataFromMap(uid)
	} else {
		// 不是本机，充当代理
		return GetDataFromOtherMap(hostRequest, req)
	}
}

// 获取本机map
func (a *AccessControl) GetDataFromMap(uid int) (isOK bool) {
	data := a.GetNewRecord(uid)
	if data == nil {
		isOK = false
	}
	return true
}

// 获取其他节点的map处理结果
func GetDataFromOtherMap(host string, request *http.Request) (isOK bool) {
	hostUrL := "http://" + host + ":" + port + "/checkRight"
	resp, body, err := GetCurl(hostUrL, request)
	if err != nil {
		golog.Error(err)
		return false
	}
	if resp.StatusCode == 200 && string(body) == "true" {
		return true
	}
	return false
}

// 模拟请求访问
func GetCurl(hostUrl string, request *http.Request) (resp *http.Response, body []byte, err error) {
	uidCookie, err := request.Cookie("uid")
	if err != nil {
		golog.Error("获取分布式锁时，获取uid失败")
		return nil, nil, err
	}
	signCookie, err := request.Cookie("sign")
	if err != nil {
		golog.Error("获取分布式锁时，获取sign失败", err)
		return nil, nil, err
	}
	client := http.DefaultClient
	req, err := http.NewRequest("GET", hostUrl, nil)
	if err != nil {
		golog.Error("创建http请求出错", err)
		return
	}
	// 将cookie注入请求
	req.AddCookie(uidCookie)
	req.AddCookie(signCookie)
	// 执行请求动作，并获取响应
	resp, err = client.Do(req)
	if err != nil {
		golog.Error("发送http请求出错", err)
		return nil, nil, err
	}
	defer resp.Body.Close()
	//
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		golog.Error("解析http响应出错", err)
		return
	}
	return
}

func main() {
	hashConsistent = common.NewConsistent()
	// 服务器添加到hash环上
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}

	//localIP, err := common.GetIntranceIP()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//localHost = localIP
	fmt.Println(localHost)
	rabbitmqValidate = rabbitmq.NewRabbitMQSimple("imoocProduct")

	filter := common.NewFilter()
	filter.RegisterFilterUri("/check", Auth)
	filter.RegisterFilterUri("/checkRight", Auth)
	http.HandleFunc("/check", filter.Handle(Check))
	http.HandleFunc("/checkRight", filter.Handle(CheckRight))
	// 启动服务端口
	http.ListenAndServe(":"+port, nil)
}
