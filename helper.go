package goetl_starter

import (
	"sync"

	logger "github.com/kordar/gologger"

	"github.com/kordar/goetl/engine"
)

var (
	engines = make(map[string]*engine.Engine)
	mu      sync.RWMutex
)

func Get(name string) *engine.Engine {
	mu.RLock()
	defer mu.RUnlock()
	eng, ok := engines[name]
	if !ok {
		logger.Fatalf("etl engine %s not exist.", name)
	}
	return eng
}

func Provide(name string, eng *engine.Engine) {
	mu.Lock()
	defer mu.Unlock()
	engines[name] = eng
}
