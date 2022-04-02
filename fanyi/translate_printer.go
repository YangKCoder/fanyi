package fanyi

import "fmt"

type TranslatePrinter interface {
	Translate
	Printer
}

func TranslatePrint(t TranslatePrinter, queryString string) {
	data, err := t.Translate(queryString)
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Print(data)
}
