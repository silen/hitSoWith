package convert

// @Time : 2019/08/25
// @Author : silen
// @File : utils
// @Software: vscode
// @Desc: .....
import (
	"encoding/json"
	"hitSoWith/libraries/utils"

	"github.com/astaxie/beego/logs"
	"github.com/fatih/structs"
	"github.com/jinzhu/copier"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

//StringToJSON string 转json
func StringToJSON(str string) (m map[string]interface{}) {
	m = map[string]interface{}{}
	if str != "" {
		if err := json.Unmarshal([]byte(str), &m); err != nil {
			logs.Error(err, "string to json failed!----"+utils.TrimLinefeed(str),
				utils.GetlastStepFileInfoForLogs(3),
				utils.GetlastStepFileInfoForLogs(4))
		}
	}
	return
}

//JSONStringToStruct json串转成相应的struct，常用于封装装载http接口请求返回数据
// structObject要传struct类型的指针（居于性能考虑，大数据结构传递都要传指针）
func JSONStringToStruct(jsonString string, structObject interface{}) error {
	return MapToStruct(StringToJSON(jsonString), structObject)
}

//MapToStruct map[string]interface 转 相应的struct
func MapToStruct(mapdata map[string]interface{}, structObject interface{}) (err error) {
	//将 map 转换为指定的结构体
	if err = mapstructure.Decode(mapdata, structObject); err != nil {
		logs.Error(mapdata, "--map[string]interface to struct failed--", utils.TrimLinefeed(err.Error()),
			utils.GetlastStepFileInfoForLogs(3),
			utils.GetlastStepFileInfoForLogs(4),
			utils.GetlastStepFileInfoForLogs(5))
	}
	return
}

//StructToMap struct 转 map[string]interface{}
func StructToMap(s interface{}) map[string]interface{} {
	return structs.Map(s)
}

//CacheDataToStruct  缓存里的数据装载进struct
func CacheDataToStruct(data interface{}, structObject interface{}) error {
	if v, ok := data.(map[string]interface{}); ok {
		return MapToStruct(v, structObject)
	}

	logs.Error(data, "--cachedata convert to map[string]interface faild",
		utils.GetlastStepFileInfoForLogs(3),
		utils.GetlastStepFileInfoForLogs(4))
	return errors.New("convert to map[string]interface faild")
}

//StructToStruct 从目标struct装载数据进另一个struct
func StructToStruct(from, to interface{}) {
	copier.Copy(to, from)
}
