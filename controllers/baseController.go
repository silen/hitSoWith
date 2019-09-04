// @Time : 2019/08/25
// @Author : silen
// @File : baseControler
// @Software: GoLand
// @Desc: controller基类

package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/silen/hitSoWith/libraries/conf"
)

//BaseController controller 基类
type BaseController struct {
	beego.Controller
	PageIndex             string
	PageSize              string
	repeatCheckLoginCount int
}

//Prepare Method 方法之前执行
func (c *BaseController) Prepare() {

	c.PrepareTrace()
	//统一检查登录
	c.CheckLogin()

	c.InitPageParames()
}

//Finish Method 方法之后执行的
func (c *BaseController) Finish() {
	c.FinishTrace()
}

//PrintConf 打印conf
func (c *BaseController) PrintConf(key string) {
	logs.Info(conf.Config.Get(key))
}

//InitPageParames 初始化page的一些值
func (c *BaseController) InitPageParames() {
	//统一征用
	c.PageIndex = c.GetString("pageIndex", "1")
	c.PageSize = c.GetString("pageSize", "10")
}

//PrepareTrace ...
func (c *BaseController) PrepareTrace() {

}

//FinishTrace ...
func (c *BaseController) FinishTrace() {

}

//CheckLogin 检查登录
func (c *BaseController) CheckLogin() {

}

//Message 消息输出
func (c *BaseController) Message(code int, message string) {
	c.ReturnJSON(map[string]interface{}{
		"status": code,
		"data":   make([]string, 0),
	}, message)
}

//ReturnJSON 指定格式输出
func (c *BaseController) ReturnJSON(data interface{}, message ...string) {

	ret := make(map[string]interface{})

	switch data.(type) {
	case int, string, []interface{}:
		ret["status"] = 200
		ret["data"] = data

	case map[string]interface{}:
		tdata := data.(map[string]interface{})
		if _, ok := tdata["status"]; !ok {
			ret["status"] = 200
		} else {
			ret["status"] = tdata["status"]
			delete(tdata, "status")
		}
		if _, ok := tdata["data"]; !ok {
			ret["data"] = tdata
		} else {
			ret["data"] = tdata["data"]
			delete(tdata, "data")
		}
		if _, ok := tdata["message"]; !ok {
			ret["message"] = ""
		} else {
			ret["message"] = tdata["message"]
			delete(tdata, "message")
		}

	default:
		ret["status"] = 200
		ret["data"] = data
	}

	var (
		hasIndent   = true
		hasEncoding = false
	)

	if beego.BConfig.RunMode == beego.PROD {
		hasIndent = false
	}
	if len(message) > 0 && message[0] != "" {
		ret["message"] = message[0]
	}

	//===========================zipkin end!
	//c.Ctx.ResponseWriter.Header().Set("go-server-ip", utils.GetServerIP())
	//c.Ctx.ResponseWriter.Header().Set("go-start-up-time", utils.GetGoStartUpTime())

	c.Ctx.Output.JSON(ret, hasIndent, hasEncoding)
	c.Finish()
	//中断controller逻辑
	c.StopRun()
}
