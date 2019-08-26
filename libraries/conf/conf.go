// @Time : 2019/08/25
// @Author : silen
// @File : conf
// @Software: vsCode
// @Desc: to do somewhat..

package conf

import (
	"os"

	"github.com/astaxie/beego"
)

//Conf conf类型
type Conf map[string]string

//Config 初始化
var Config = Conf{
	"logpath":                      "/apps/logs/golang/",
	"logfile":                      "",
	"logconsole":                   "true",
	"logLevel":                     "7",
	"HttpTimeout":                  "20",
	"HttpErrRetryTimes":            "0",
	"HystrixTimeout":               "20",
	"HystrixMaxConcurrentRequests": "1000",
	"HystrixErrorPercentThreshold": "50",

	"mysql_host":     "127.0.0.1",
	"mysql_port":     "3306",
	"mysql_username": "isp_test",
	"mysql_password": "nDvsEQr4cvk5",
	"mysql_database": "isp_test",
	"mysql_charset":  "utf8",
	"timezone":       "Asia/Shanghai",
	"redis_host":     "127.0.0.1",
	"redis_port":     "6379",
	"redis_password": "8xSFDJJnVxuv",
}

//Get 根据key获取string值
func (c *Conf) Get(key string) (v string) {
	//优先从本应用的conf拿，优先支持私有配置
	v = beego.AppConfig.String(key)
	//找不到再去系统环境变量尝试找
	if v == "" {
		v = os.Getenv(key)
	}

	//再不行就拿这里配的默认值（默认配置）
	if v == "" {
		if vv, exists := (*c)[key]; exists {
			v = vv
		}
	}
	return
}
