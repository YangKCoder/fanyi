package fanyi

import "fmt"

type Translate interface {
	Translate(queryString string) ([]byte, error)
}

type Printer interface {
	Print(data []byte)
}

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
