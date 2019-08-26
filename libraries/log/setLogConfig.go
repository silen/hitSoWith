// @Time : 2019/3/13 4:23 PM
// @Author : acol
// @File : setLogConfig
// @Software: GoLand
// @Desc: 统一设置项目的日志文件配置

package log

import (
	"fmt"
	"os"
	"time"

	"github.com/astaxie/beego"
	"gitlab.sys.hxsapp.net/isp/isp-common-go.git/libraries/conf"
	"gitlab.sys.hxsapp.net/isp/isp-common-go.git/libraries/utils"
)

func init() {
	//在服务器上每次启动都为守护程序systemd记录一下pid
	if beego.BConfig.RunMode != "dev" {
		//writePidInFile()
	}

	//期望输出调用的文件名和文件行号
	beego.BeeLogger.EnableFuncCallDepth(true)
	//文件名和行号获取深度
	beego.BeeLogger.SetLogFuncCallDepth(3)

	//是否关闭控制台输出
	logconsole := conf.Config.Get("logconsole")
	if logconsole == "false" {
		beego.BeeLogger.DelLogger("console")
	}

	logfile := conf.Config.Get("logfile")
	//如果日志文件名设置为false，则日志不写文件
	if logfile == "false" {
		return
	}
	if logfile == "" {
		logfile = conf.Config.Get("appname") + ".log"
	}

	var filepath = "."
	//设置日志文件路径，如果logpath为.则日志文件默认放在当前程序根目录下
	if conf.Config.Get("logpath") != "." {
		filepath1 := conf.Config.Get("logpath")
		makePathExists(filepath1)

		filepath = filepath1 + conf.Config.Get("appname")
		makePathExists(filepath)
	}

	//设置日志文件，文件名带日期格式
	filename := fmt.Sprintf("%s/%s", filepath, logfile)
	beego.BeeLogger.SetLogger("file", `{"filename":"`+filename+`", "maxdays":30}`)
	beego.BeeLogger.SetLevel(conf.Config.GetInt("logLevel"))

	//设置日志为异步输出，并设置内存缓存chan大小为1000
	beego.BeeLogger.Async(1e3)

	//sleep 50毫秒，确保以上设置成功
	time.Sleep(50 * time.Millisecond)
	beego.BeeLogger.Info("===BEEGO_ENV===" + os.Getenv("BEEGO_MODE") + "===" + os.Getenv("BEEGO_RUNMODE"))
	beego.BeeLogger.Info("====BEEGO_RUNMODE===" + beego.BConfig.RunMode + "====" + conf.Config.Get("runmode"))
	beego.BeeLogger.Info("===============set log config success!==========")
	beego.BeeLogger.Info("=========filename:" + filename + "=========")
	utils.SetGoStartUpTime()
}

func makePathExists(filepath string) {
	_, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			beego.BeeLogger.Info(filepath + " is not exist")
			err = os.Mkdir(filepath, os.ModePerm)
			if err != nil {
				beego.BeeLogger.Error(filepath + " mkdir failed!")
			}
		}
	}
}
