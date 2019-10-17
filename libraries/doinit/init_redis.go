package doinit

import (
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/silen/hitSoWith/libraries/cache"
	"github.com/silen/hitSoWith/libraries/conf"

	//导入框架定制的redis
	_ "github.com/silen/hitSoWith/libraries/cache/redis"
)

var (
	//RDS rds对象
	RDS cache.Cache
)

//InitRedis 初始化redis
func InitRedis() {
	if RDS != nil {
		//logs.Info("redis already init", RDS)
		return
	}
	logs.Info("redis init...")

	host := conf.Config.Get("redis_host")
	port := conf.Config.Get("redis_port")
	password := conf.Config.Get("redis_password")

	connStr := fmt.Sprintf("{\"conn\":\"%s:%s\",\"password\":\"%s\"}", host, port, password)

	bm, err := cache.NewCache("redis", connStr)

	if err == nil {
		RDS = bm
	} else {
		logs.Error("redis init error====", err)
	}
}
