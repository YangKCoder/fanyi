package factory

import (
	"fanyi/fanyi"
	"fmt"
	"sync"
)

var (
	providerMu sync.RWMutex
	providers  = make(map[string]fanyi.TranslatePrinter)
)

func Register(name string, p fanyi.TranslatePrinter) {
	providerMu.Lock()
	defer providerMu.Unlock()

	if p == nil {
		panic("store: Register provider is nil")
	}
	if _, dup := providers[name]; dup {
		panic("store: Register called twice for provider " + name)
	}
	providers[name] = p
}

func New(providerName string) (fanyi.TranslatePrinter, error) {
	providerMu.RLock()
	p, ok := providers[providerName]
	if !ok {
		return nil, fmt.Errorf("store: unkonwn provider %s", providerName)
	}
	return p, nil
}

func Providers() map[string]fanyi.TranslatePrinter {
	return providers
}
