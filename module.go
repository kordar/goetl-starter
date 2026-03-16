package goetl_starter

import (
	"fmt"

	logger "github.com/kordar/gologger"

	"github.com/kordar/goetl/engine"
	"github.com/spf13/cast"
)

type EngineLoaderR func(moduleName string, itemID string, item map[string]any) (*engine.Engine, error)
type EngineLoader func(moduleName string, itemID string, item map[string]any)

type GoetlModule struct {
	name string
	load any
}

func NewGoetlModule(name string, load any) *GoetlModule {
	return &GoetlModule{name: name, load: load}
}

func (m GoetlModule) Name() string {
	return m.name
}

func (m GoetlModule) _load(id string, cfg map[string]any) {
	if id == "" {
		logger.Fatalf("[%s] the attribute id cannot be empty.", m.Name())
		return
	}
	if m.load == nil {
		logger.Fatalf("[%s] load callback cannot be nil. id=%s", m.Name(), id)
		return
	}

	eng, err, provided := m.callLoad(id, cfg)
	if err != nil {
		logger.Fatalf("[%s] id=%s err=%v", m.Name(), id, err)
		return
	}

	if provided {
		logger.Infof("[%s] loading module '%s' successfully", m.Name(), id)
		return
	}

	if eng == nil {
		logger.Warnf("[%s] id=%s engine is nil", m.Name(), id)
		return
	}

	Provide(id, eng)
	logger.Infof("[%s] loading module '%s' successfully", m.Name(), id)
}

func (m GoetlModule) Load(value any) {
	if value == nil {
		return
	}

	items := cast.ToStringMap(value)
	if items["id"] != nil {
		id := cast.ToString(items["id"])
		m._load(id, items)
		return
	}

	for key, item := range items {
		m._load(key, cast.ToStringMap(item))
	}
}

func (m GoetlModule) Close() {
}

func (m GoetlModule) callLoad(id string, cfg map[string]any) (*engine.Engine, error, bool) {
	switch f := m.load.(type) {
	case EngineLoaderR:
		eng, err := f(m.Name(), id, cfg)
		return eng, err, false
	case func(moduleName string, itemID string, item map[string]any) (*engine.Engine, error):
		eng, err := f(m.Name(), id, cfg)
		return eng, err, false
	case func(moduleName string, itemID string, item map[string]any) *engine.Engine:
		eng := f(m.Name(), id, cfg)
		return eng, nil, false
	case func(itemID string, item map[string]any) (*engine.Engine, error):
		eng, err := f(id, cfg)
		return eng, err, false
	case func(itemID string, item map[string]any) *engine.Engine:
		eng := f(id, cfg)
		return eng, nil, false
	case func(item map[string]any) (*engine.Engine, error):
		eng, err := f(cfg)
		return eng, err, false
	case func(item map[string]any) *engine.Engine:
		eng := f(cfg)
		return eng, nil, false
	case EngineLoader:
		f(m.Name(), id, cfg)
		return getProvided(id), nil, getProvided(id) != nil
	case func(moduleName string, itemID string, item map[string]any):
		f(m.Name(), id, cfg)
		return getProvided(id), nil, getProvided(id) != nil
	case func(itemID string, item map[string]any):
		f(id, cfg)
		return getProvided(id), nil, getProvided(id) != nil
	case func(item map[string]any):
		f(cfg)
		return getProvided(id), nil, getProvided(id) != nil
	default:
		return nil, fmt.Errorf("unsupported load callback type: %T", m.load), false
	}
}

func getProvided(id string) *engine.Engine {
	mu.RLock()
	eng := engines[id]
	mu.RUnlock()
	return eng
}
