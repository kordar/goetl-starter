package goetlstarter

import (
	"fmt"
	"log/slog"

	"github.com/kordar/goetl/engine"
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
		slog.Error("the attribute id cannot be empty", "module", m.Name())
		panic(fmt.Errorf("[%s] id empty", m.Name()))
	}
	if m.load == nil {
		slog.Error("load callback cannot be nil", "module", m.Name(), "id", id)
		panic(fmt.Errorf("[%s] load callback nil", m.Name()))
	}

	eng, err, provided := m.callLoad(id, cfg)
	if err != nil {
		slog.Error("load error", "module", m.Name(), "id", id, "err", err)
		panic(err)
	}

	if provided {
		slog.Info("loading module successfully", "module", m.Name(), "id", id)
		return
	}

	if eng == nil {
		slog.Warn("engine is nil", "module", m.Name(), "id", id)
		return
	}

	Provide(id, eng)
	slog.Info("loading module successfully", "module", m.Name(), "id", id)
}

func (m GoetlModule) Load(value any) {
	if value == nil {
		return
	}

	items := toStringMap(value)
	if items["id"] != nil {
		id := toString(items["id"])
		m._load(id, items)
		return
	}

	for key, item := range items {
		m._load(key, toStringMap(item))
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

func toString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case fmt.Stringer:
		return x.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func toStringMap(v any) map[string]any {
	if v == nil {
		return map[string]any{}
	}
	switch m := v.(type) {
	case map[string]any:
		return m
	case map[any]any:
		out := make(map[string]any, len(m))
		for k, val := range m {
			out[toString(k)] = val
		}
		return out
	default:
		return map[string]any{}
	}
}
