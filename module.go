package goetl_starter

import (
	"log"

	"github.com/kordar/goetl/engine"
	"github.com/spf13/cast"
)

type EngineLoader func(moduleName string, itemID string, item map[string]any) (*engine.Engine, error)

type GoetlModule struct {
	name string
	load EngineLoader
}

func NewGoetlModule(name string, load EngineLoader) *GoetlModule {
	return &GoetlModule{name: name, load: load}
}

func (m GoetlModule) Name() string {
	return m.name
}

func (m GoetlModule) _load(id string, cfg map[string]any) {
	if id == "" {
		log.Fatalf("[%s] the attribute id cannot be empty.", m.Name())
		return
	}
	if m.load == nil {
		log.Fatalf("[%s] load callback cannot be nil. id=%s", m.Name(), id)
		return
	}

	eng, err := m.load(m.Name(), id, cfg)
	if err != nil {
		log.Fatalf("[%s] id=%s err=%v", m.Name(), id, err)
		return
	}
	if eng == nil {
		log.Fatalf("[%s] id=%s engine is nil", m.Name(), id)
		return
	}

	Provide(id, eng)
	log.Printf("[%s] loading module '%s' successfully", m.Name(), id)
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
