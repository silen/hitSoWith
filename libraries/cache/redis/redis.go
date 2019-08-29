// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package redis for cache provider
//
// depend on github.com/gomodule/redigo/redis
//
// go install github.com/gomodule/redigo/redis
//
// Usage:
// import(
//   _ "github.com/astaxie/beego/cache/redis"
//   "github.com/astaxie/beego/cache"
// )
//
//  bm, err := cache.NewCache("redis", `{"conn":"127.0.0.1:11211"}`)
//
//  more docs http://beego.me/docs/module/cache.md
package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/silen/hitSoWith/libraries/cache"

	"strings"
)

var (
	// DefaultKey the collection name of redis for cache adapter.
	DefaultKey = "xyRDS"
)

// Cache is Redis cache adapter.
type Cache struct {
	p        *redis.Pool // redis connection pool
	conninfo string
	dbNum    int
	key      string
	password string
	maxIdle  int
}

// NewRedisCache create new redis cache with default collection name.
func NewRedisCache() cache.Cache {
	return &Cache{key: DefaultKey}
}

// actually do the redis cmds, args[0] must be the key name.
func (rc *Cache) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	args[0] = rc.associate(args[0])
	c := rc.p.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

// associate with config key.
func (rc *Cache) associate(originKey interface{}) string {
	//这他妈的会在所有key面前拼前缀！！
	//这样其它栈如果共同操作就要适配这个前缀，算了。。。注释掉
	//acol	2019-04-16
	//return fmt.Sprintf("%s:%s", rc.key, originKey)
	return fmt.Sprintf("%s", originKey)
}

// Get cache from redis.
func (rc *Cache) Get(key string) interface{} {
	if v, err := rc.do("GET", key); err == nil {
		return v
	}
	return nil
}

// MGet get cache from redis.
func (rc *Cache) MGet(keys []string) []interface{} {
	c := rc.p.Get()
	defer c.Close()
	var args []interface{}
	for _, key := range keys {
		args = append(args, rc.associate(key))
	}
	values, err := redis.Values(c.Do("MGET", args...))
	if err != nil {
		return nil
	}
	return values
}

// Set put cache to redis.
func (rc *Cache) Set(key string, val interface{}, timeout time.Duration) error {
	var err error
	if timeout > time.Duration(0) {
		_, err = rc.do("SETEX", key, int64(timeout/time.Second), val)
	} else {
		_, err = rc.do("SET", key, val)
	}

	return err
}

//==============Hash================

// HGet cache from redis.
func (rc *Cache) HGet(key string, field string) interface{} {

	if v, err := rc.do("HGET", key, field); err == nil {

		return v

	}

	return nil

}

//HMGet hmget cache
func (rc *Cache) HMGet(key string, fields []string) []interface{} {
	c := rc.p.Get()
	defer c.Close()
	var args []interface{}
	args = append(args, rc.associate(key))
	for _, field := range fields {
		args = append(args, field)
	}

	values, err := redis.Values(c.Do("HMGET", args...))

	if err != nil {
		return nil
	}

	return values
}

// HGetAll cache from redis.
func (rc *Cache) HGetAll(key string) []interface{} {
	ret := make([]interface{}, 1)
	if v, err := redis.Values(rc.do("HGETALL", key)); err == nil {
		for _, value := range v {
			//fmt.Printf("%s ", value.([]byte))
			ret = append(ret, value.([]byte))
		}
	}
	return ret
}

//HSet set hash cache
func (rc *Cache) HSet(key string, field string, val interface{}) error {

	var err error

	_, err = rc.do("HSET", key, field, val)

	return err

}

//HLen len hash cache
func (rc *Cache) HLen(key string) int {
	v, err := redis.Int(rc.do("HLEN", key))
	if err != nil {
		return 0
	}
	return v
}

//HExists check hash cache exists
func (rc *Cache) HExists(key string, field string) bool {
	_, err := redis.Int(rc.do("HEXISTS", key, field))
	if err != nil {
		return false
	}
	return true
}

//HDel len hash cache
func (rc *Cache) HDel(key string, fields []string) int {
	c := rc.p.Get()
	defer c.Close()
	var args []interface{}
	args = append(args, rc.associate(key))
	for _, field := range fields {
		args = append(args, field)
	}

	values, err := redis.Int(c.Do("HDEL", args...))

	if err != nil {
		return 0
	}

	return values
}

//==============Hash end============

//==============List================

//LPush lpush list cache
//将一个或多个值value插入到列表key的表头
//返回列表长度
func (rc *Cache) LPush(key string, val ...interface{}) int {
	var args []interface{}
	args = append(args, key)
	for _, v := range val {
		args = append(args, v)
	}
	v, err := rc.do("LPUSH", args...)
	l, _ := redis.Int(v, err)
	return l
}

//RPush rpush list cache
//将一个或多个值value插入到列表key的表尾
//返回列表长度
func (rc *Cache) RPush(key string, val ...interface{}) int {
	var args []interface{}
	args = append(args, key)
	for _, v := range val {
		args = append(args, v)
	}
	v, err := rc.do("RPUSH", args...)
	l, _ := redis.Int(v, err)
	return l
}

//LPop lpop list cache
//移除并返回列表key的头元素
func (rc *Cache) LPop(key string) interface{} {
	v, err := rc.do("LPOP", key)
	if err != nil {
		return nil
	}
	return v
}

//RPop lpop list cache
//移除并返回列表key的尾元素
func (rc *Cache) RPop(key string) interface{} {
	v, err := rc.do("RPOP", key)
	if err != nil {
		return nil
	}
	return v
}

//LLen llen list cache
//返回列表key的长度。
func (rc *Cache) LLen(key string) int {
	v, err := rc.do("LLEN", key)
	l, _ := redis.Int(v, err)
	return l
}

//LRange lrange list cache
//返回列表key中指定区间内的元素，区间以偏移量start和stop指定。
func (rc *Cache) LRange(key string, start int, stop int) []interface{} {
	values, err := redis.Values(rc.do("LRANGE", key, start, stop))
	if err != nil {
		return nil
	}

	return values
}

//LRem lrem list cache
//根据参数count的值，移除列表中与参数value相等的元素。
//
//count的值可以是以下几种：
//count > 0: 从表头开始向表尾搜索，移除与value相等的元素，数量为count。
//count < 0: 从表尾开始向表头搜索，移除与value相等的元素，数量为count的绝对值。
//count = 0: 移除表中所有与value相等的值。
//返回被移除的元素数量
func (rc *Cache) LRem(key string, count int, value interface{}) int {
	v, err := rc.do("LREM", key, count, value)
	l, _ := redis.Int(v, err)
	return l
}

//LSet 将列表key下标为index的元素的值甚至为value
//返回操作是否成功
func (rc *Cache) LSet(key string, index int, value interface{}) bool {
	_, err := rc.do("LSET", key, index, value)
	if err != nil {
		return false
	}
	return true
}

//LIndex 返回列表key中，下标为index的元素
func (rc *Cache) LIndex(key string, index int) interface{} {
	v, err := rc.do("LINDEX", key, index)
	if err != nil {
		return nil
	}
	return v
}

//LTrim 对一个列表进行修剪(trim)，就是说，让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除。
//举个例子，执行命令LTRIM list 0 2，表示只保留列表list的前三个元素，其余元素全部删除。
//返回操作是否成功
func (rc *Cache) LTrim(key string, start int, stop int) bool {
	_, err := rc.do("LTRIM", key, start, stop)
	if err != nil {
		return false
	}
	return true
}

//===========List end================

// Delete delete cache in redis.
func (rc *Cache) Delete(key string) error {
	_, err := rc.do("DEL", key)
	return err
}

// IsExist check cache's existence in redis.
func (rc *Cache) IsExist(key string) bool {
	v, err := redis.Bool(rc.do("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

//SetExpire set timeout
func (rc *Cache) SetExpire(key string, timeout time.Duration) error {
	var err error
	_, err = rc.do("EXPIRE", key, int64(timeout/time.Second))
	return err
}

//TTL get cache's ttl
func (rc *Cache) TTL(key string) int64 {
	v, err := redis.Int64(rc.do("TTL", key))

	if err != nil {
		return 0
	}

	return v
}

//Keys get all keys
func (rc *Cache) Keys(key string) []string {
	v, err := rc.do("KEYS", key)

	vv, err2 := redis.Strings(v, err)
	if err2 != nil {
		vv = make([]string, 0)
	}

	return vv
}

// Incr increase counter in redis.
func (rc *Cache) Incr(key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, 1))
	return err
}

// Decr decrease counter in redis.
func (rc *Cache) Decr(key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, -1))
	return err
}

// ClearAll clean all cache in redis. delete this redis collection.
func (rc *Cache) ClearAll() error {
	c := rc.p.Get()
	defer c.Close()
	cachedKeys, err := redis.Strings(c.Do("KEYS", rc.key+":*"))
	if err != nil {
		return err
	}
	for _, str := range cachedKeys {
		if _, err = c.Do("DEL", str); err != nil {
			return err
		}
	}
	return err
}

// StartAndGC start redis cache adapter.
// config is like {"key":"collection key","conn":"connection info","dbNum":"0"}
// the cache item in redis are stored forever,
// so no gc operation.
func (rc *Cache) StartAndGC(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)

	if _, ok := cf["key"]; !ok {
		cf["key"] = DefaultKey
	}
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}

	// Format redis://<password>@<host>:<port>
	cf["conn"] = strings.Replace(cf["conn"], "redis://", "", 1)
	if i := strings.Index(cf["conn"], "@"); i > -1 {
		cf["password"] = cf["conn"][0:i]
		cf["conn"] = cf["conn"][i+1:]
	}

	if _, ok := cf["dbNum"]; !ok {
		cf["dbNum"] = "0"
	}
	if _, ok := cf["password"]; !ok {
		cf["password"] = ""
	}
	if _, ok := cf["maxIdle"]; !ok {
		cf["maxIdle"] = "3"
	}
	rc.key = cf["key"]
	rc.conninfo = cf["conn"]
	rc.dbNum, _ = strconv.Atoi(cf["dbNum"])
	rc.password = cf["password"]
	rc.maxIdle, _ = strconv.Atoi(cf["maxIdle"])

	rc.connectInit()

	c := rc.p.Get()
	defer c.Close()

	return c.Err()
}

// connect to redis.
func (rc *Cache) connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", rc.conninfo)
		if err != nil {
			return nil, err
		}

		if rc.password != "" {
			if _, err := c.Do("AUTH", rc.password); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", rc.dbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	rc.p = &redis.Pool{
		MaxIdle:     rc.maxIdle,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
}

func init() {
	cache.Register("redis", NewRedisCache)
}
