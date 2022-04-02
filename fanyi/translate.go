package fanyi

type Translate interface {
	Translate(queryString string) ([]byte, error)
}
