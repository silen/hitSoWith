package main

import (
	_ "github.com/silen/hitSoWith/libraries/log"

	"github.com/astaxie/beego"
	_ "github.com/silen/hitSoWith/routers"
)

func main() {
	beego.SetLogger("file", `{"filename":"`+beego.AppConfig.String("logfile")+`"}`)
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.SetLogFuncCall(true)
	beego.Run()
}
