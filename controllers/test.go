package controllers

// @Time : 2019/08/25
// @Author : silen
// @File : test
// @Software: vsCode
// @Desc: test...
import (
	"github.com/astaxie/beego/logs"
)

//Test test
type Test struct {
	BaseController
}

//So So
func (c *Test) So() {
	var i = 1000000000000000000 //int最大18个0
	logs.Alert("So！, hit So With", i)
	c.ReturnJSON("唯一纯白的茉莉花～ biu ")
}
