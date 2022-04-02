package cmd

import (
	"fanyi/fanyi"
	"fanyi/fanyi/factory"
	_ "fanyi/internal"
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"net/url"
	"os"
	"strings"
	"sync"
)

func Execute() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s word\n\n", os.Args[0])
		flag.PrintDefaults()
		eg := `Examples:
  $ fanyi word
  $ fanyi world peace
  $ fanyi 中文`
		fmt.Println(eg)
	}
	flag.Parse()
	var queryString string
	if len(os.Args[1:]) == 0 {
		text, err := clipboard.ReadAll()
		if err != nil || text == "" {
			//读取剪切板失败或者没内容
			flag.Usage()
			return
		}
		fmt.Printf(" \n 默认读取剪贴板: %s\n", text)
		queryString = text
	} else {
		queryString = strings.Join(flag.Args(), " ")
	}
	queryString = url.QueryEscape(queryString)
	fmt.Println()
	providers := factory.Providers()
	group := sync.WaitGroup{}
	for _, p := range providers {
		group.Add(1)
		go func(f fanyi.TranslatePrinter) {
			fanyi.TranslatePrint(f, queryString)
			group.Done()
		}(p)
	}
	group.Wait()
}
