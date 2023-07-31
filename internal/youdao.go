package internal

import (
	"fanyi/fanyi/factory"
	"fanyi/internal/utils"
	"fanyi/internal/utils/authv3"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
)

const youdaoUrl = "https://openapi.youdao.com/api"

var appKey = os.Getenv("YOUDAO_APPKEY")
var appSecret = os.Getenv("YOUDAO_APPSECRET")

type YouDao struct{}

func init() {
	if appKey == "" || appSecret == "" {
		return
	}
	factory.Register("youdao", &YouDao{})
}

func (y *YouDao) Translate(queryString string) ([]byte, error) {

	// 添加请求参数
	paramsMap := map[string][]string{
		"q":    {queryString},
		"from": {"auto"},
		"to":   {"en"},
	}
	header := map[string][]string{
		"Content-Type": {"application/x-www-form-urlencoded"},
	}
	// 添加鉴权相关参数
	authv3.AddAuthParams(appKey, appSecret, paramsMap)
	// 请求api服务
	result := utils.DoPost(youdaoUrl, header, paramsMap, "application/json")

	return result, nil
}

func (y *YouDao) Print(data []byte) {
	json, err := simplejson.NewJson(data)
	if err != nil {
		log.Fatal(err)
	}
	query, _ := json.Get("query").String()
	phonetic, err := json.Get("basic").Get("phonetic").String()
	var phoneticStr string
	if err != nil {
		phoneticStr = ""
	} else {
		phoneticStr = fmt.Sprintf("[ %s ]", magenta(phonetic))
	}
	fmt.Printf(" %s %s %s\n\n", query, phoneticStr, gray("~  fanyi.youdao.com"))
	explains, _ := json.Get("basic").Get("explains").Array()
	for _, value := range explains {
		fmt.Printf(" %s %s\n", gray("-"), green(value))
	}
	fmt.Println()
	web, _ := json.Get("web").Array()
	for i, value := range web {
		val := value.(map[string]interface{})
		fmt.Printf(" %s %s\n", gray(strconv.Itoa(i+1)+"."), highlight(val["key"].(string), query))
		valuelen := len(val["value"].([]interface{}))
		valArr := make([]string, valuelen)
		for i, value := range val["value"].([]interface{}) {
			valArr[i] = value.(string)
		}
		valueStr := strings.Join(valArr, ", ")
		fmt.Printf("    %s\n", cyan(valueStr))
	}
	fmt.Println()
	fmt.Println(gray("   --------"))
	fmt.Println()
}
