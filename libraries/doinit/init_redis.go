package doinit

import (
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/silen/hitSoWith/libraries/cache"
	"github.com/silen/hitSoWith/libraries/conf"

	//vendor redis
	_ "github.com/silen/hitSoWith/libraries/cache/redis"
)

var (
	//RDS rds对象
	RDS cache.Cache
)

//NewRedis 初始化redis
func NewRedis() {
	if RDS != nil {
		//logs.Info("redis already init", RDS)
		return
	}

	host := conf.Config.Get("redis_host")
	port := conf.Config.Get("redis_port")
	password := conf.Config.Get("redis_password")
	bm, err := cache.NewCache("redis", fmt.Sprintf("{\"conn\":\"%s:%s\",\"password\":\"%s\"}", host, port, password))
	if err == nil {
		RDS = bm
	} else {
		logs.Error("redis init error====", err)
	}
}
