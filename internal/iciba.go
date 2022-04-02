package internal

import (
	"encoding/xml"
	"errors"
	"fanyi/fanyi/factory"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const icibaUrl = "https://dict-co.iciba.com/api/dictionary.php?key=D191EBD014295E913574E1EAF8E06666&w=${word}"

func init() {
	factory.Register("iciba", &Iciba{})
}

type Iciba struct { //xmlName     xml.Name `xml:"dict"`
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

func (v *Iciba) Translate(queryString string) ([]byte, error) {
	url := strings.Replace(icibaUrl, "${word}", queryString, 1)
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("iciba翻译接口问题")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return data, nil
}

func (v *Iciba) Print(data []byte) {
	err := xml.Unmarshal(data, v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	var phoneticStr string
	for i, value := range v.Ps {
		if i == 0 {
			phoneticStr += "英" + "[ " + value + "] "
		} else {
			phoneticStr += "美" + "[ " + value + "] "
		}
	}
	fmt.Printf(" %s %s %s\n\n", v.Key, magenta(phoneticStr), gray("~  iciba.com"))
	if !isChinese(v.Key) {
		for i := 0; i < len(v.Pos); i++ {
			fmt.Printf(" %s %s %s", gray("-"), green(v.Pos[i]), green(v.Acceptation[i]))
		}
	}
	fmt.Println()
	for i := 0; i < len(v.Sent); i++ {
		fmt.Printf(" %s %s\n", gray(strconv.Itoa(i+1)+"."), highlight(del(v.Sent[i].Orig), v.Key))
		fmt.Printf("    %s\n", cyan(del(v.Sent[i].Trans)))
	}
	fmt.Println()
	fmt.Println(gray("   --------"))
	fmt.Println()
}

// 删除string中的换行符
func del(str string) string {
	r := regexp.MustCompile("\n")
	res := r.ReplaceAllString(str, "")
	return res
}

// 是否包含中文
func isChinese(str string) bool {
	count := 0
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			count++
		}
	}
	return count > 0
}
