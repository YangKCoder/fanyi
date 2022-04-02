package internal

import (
	"errors"
	"fanyi/fanyi/factory"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const youdaoUrl = "https://fanyi.youdao.com/openapi.do?keyfrom=node-fanyi&key=110811608&type=data&doctype=json&version=1.1&q=${word}"

type YouDao struct {
	Key         string   `xml:"key"`
	Ps          []string `xml:"ps"`
	Pron        []string `xml:"pron"`
	Pos         []string `xml:"pos"`
	Acceptation []string `xml:"acceptation"`
	Sent        []struct {
		Orig  string `xml:"orig"`
		Trans string `xml:"trans"`
	} `xml:"sent"`
}

func init() {
	factory.Register("youdao", &YouDao{})
}

func (y *YouDao) Translate(queryString string) ([]byte, error) {
	url := strings.Replace(youdaoUrl, "${word}", queryString, 1)
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("有道翻译接口问题")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return data, nil
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
