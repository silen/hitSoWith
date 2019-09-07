// @Time : 2018/12/18 5:03 PM
// @Author : acol
// @File : http
// @Software: GoLand
// @Desc: 项目http post/get 统一封装

package http

import (
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/silen/hitSoWith/libraries/conf"
)

var (
	//errRetryCount = 0 //失败重试计数
	connectTimeout = 5 //连接超时时间，固定5秒吧

	httpTimeout = conf.Config.GetInt("HttpTimeout") //http请求超时时间

)

//NewHTTP 新建Http对象
func NewHTTP() *HTTP {

	//default
	hp := &HTTP{
		connectTimeout:   time.Duration(connectTimeout) * time.Second,
		readWriteTimeout: time.Duration(httpTimeout) * time.Second, //请求超时时间
		debug:            false,
	}
	return hp
}

//Get get
func Get(url string, data ...map[string]string) string {
	hp := NewHTTP()

	tdata := make(map[string]string)
	if len(data) > 0 {
		tdata = data[0]
	}
	str, _ := hp.Get(url, tdata)

	return str
}

//Post post
func Post(url string, data ...map[string]string) string {
	hp := NewHTTP()

	tdata := make(map[string]string)
	if len(data) > 0 {
		tdata = data[0]
	}
	str, _ := hp.Post(url, tdata)

	return str
}

//HTTP 对象
type HTTP struct {
	header map[string]string
	// Body adds request raw body.
	// it supports string and []byte.
	body interface{}
	// JSONBody adds request raw body encoding by JSON.
	jsonBody         interface{}
	connectTimeout   time.Duration
	readWriteTimeout time.Duration
	debug            bool
	errRetryCount    int //失败重试计数
}

//SetHeader 设置header
func (h *HTTP) SetHeader(header map[string]string) *HTTP {
	h.header = header
	return h
}

//SetTimeout 设置超时时间,timeout单位为秒
func (h *HTTP) SetTimeout(timeout int) *HTTP {
	h.connectTimeout = time.Duration(connectTimeout) * time.Second
	h.readWriteTimeout = time.Duration(timeout) * time.Second
	return h
}

//SetDebug 是否开启debug
func (h *HTTP) SetDebug(isDebug bool) *HTTP {
	h.debug = isDebug
	return h
}

//SetBody 设置body传参
func (h *HTTP) SetBody(data interface{}) *HTTP {
	h.body = data
	return h
}

//SetJSONBody 设置body传参json格式
func (h *HTTP) SetJSONBody(data interface{}) *HTTP {
	h.jsonBody = data
	return h
}

func (h *HTTP) send(method string, url string, data map[string]string) (retStr string, retErr error) {
	if url == "" {
		return
	}

	var req *httplib.BeegoHTTPRequest
	if method == "POST" {
		req = httplib.Post(url)
	} else {
		req = httplib.Get(url)
	}

	//是否开启debug
	req.Debug(h.debug)

	//设置超时时间
	req.SetTimeout(h.connectTimeout, h.readWriteTimeout)

	//Go 语言的 HTTP 默认使用的是长连接，也就是说当请求完成之后，TCP 连接还会继续保留，直到一段时间之后（默认3分钟左右）才关闭
	//这里发出的http请求均认为是短连接，传输完成后立即关闭链接
	req.Header("Connection", "close")

	if len(h.header) > 0 {
		for k, v := range h.header {
			req.Header(k, v)
		}

		//重置header
		h.header = make(map[string]string)
	}

	//从body传参，支持string and []byte.
	if h.body != nil {
		req.Body(h.body)
	}
	//从body传参，json格式
	if h.jsonBody != nil {
		req.JSONBody(h.jsonBody)
	}

	//设置请求参数
	if len(data) > 0 {
		for k, v := range data {
			req.Param(k, v)
		}
	}

	retStr, retErr = req.String()
	return
}

//Get http get方法
func (h *HTTP) Get(url string, data map[string]string) (str string, err error) {
	return h.send("GET", url, data)
}

//Post http post 方法
func (h *HTTP) Post(url string, data map[string]string) (str string, err error) {
	return h.send("POST", url, data)
}
