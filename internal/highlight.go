package internal

import (
	"github.com/gookit/color"
	"regexp"
)

var magenta = color.FgMagenta.Render
var gray = color.FgGray.Render
var green = color.FgGreen.Render
var cyan = color.FgCyan.Render

// 高亮句子中的单词
func highlight(str string, query string) string {
	yellow := color.FgYellow.Render
	// 句子中单词用黄色，其他用灰色
	r := regexp.MustCompile("(?i)" + "(.*)" + "(" + query + ")" + "(.*)")
	res1 := r.ReplaceAllString(str, "$1$2"+gray("$3"))
	r2 := regexp.MustCompile("(?i)" + "(.*?)" + "(" + query + ")")
	res2 := r2.ReplaceAllString(res1, gray("$1")+yellow("$2"))
	return res2
}
