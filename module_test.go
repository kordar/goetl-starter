package goetlstarter

import (
	"fmt"
	"testing"
	"time"

	"github.com/kordar/goetl/engine"
)

func TestGoetlModule_Load_Multi(t *testing.T) {
	t.Parallel()

	prefix := fmt.Sprintf("test_%d_", time.Now().UnixNano())
	created := map[string]*engine.Engine{}

	m := NewGoetlModule("gorm", func(moduleName string, itemID string, item map[string]any) (*engine.Engine, error) {
		_ = moduleName
		_ = item
		eng := &engine.Engine{}
		created[itemID] = eng
		return eng, nil
	})

	m.Load(map[string]any{
		prefix + "a": map[string]any{"dsn": "a"},
		prefix + "b": map[string]any{"dsn": "b"},
	})

	if Get(prefix+"a") != created[prefix+"a"] {
		t.Fatalf("engine mismatch: a")
	}
	if Get(prefix+"b") != created[prefix+"b"] {
		t.Fatalf("engine mismatch: b")
	}
}

func TestGoetlModule_Load_SingleByID(t *testing.T) {
	t.Parallel()

	id := fmt.Sprintf("test_%d_solo", time.Now().UnixNano())
	var got *engine.Engine

	m := NewGoetlModule("gorm", func(moduleName string, itemID string, item map[string]any) (*engine.Engine, error) {
		_ = moduleName
		_ = item
		if itemID != id {
			t.Fatalf("unexpected id: %s", itemID)
		}
		got = &engine.Engine{}
		return got, nil
	})

	m.Load(map[string]any{"id": id, "dsn": "x"})

	if Get(id) != got {
		t.Fatalf("engine mismatch: solo")
	}
}

func TestGoetlModule_Load_VoidLoader_Provides(t *testing.T) {
	t.Parallel()

	id := fmt.Sprintf("test_%d_void", time.Now().UnixNano())
	want := &engine.Engine{}

	m := NewGoetlModule("gorm", func(moduleName string, itemID string, item map[string]any) {
		_ = moduleName
		_ = item
		if itemID != id {
			t.Fatalf("unexpected id: %s", itemID)
		}
		Provide(itemID, want)
	})

	m.Load(map[string]any{"id": id, "dsn": "x"})

	if Get(id) != want {
		t.Fatalf("engine mismatch: void")
	}
}

func TestGoetlModule_Load_NoErrorReturnLoader(t *testing.T) {
	t.Parallel()

	id := fmt.Sprintf("test_%d_noerr", time.Now().UnixNano())
	want := &engine.Engine{}

	m := NewGoetlModule("gorm", func(moduleName string, itemID string, item map[string]any) *engine.Engine {
		_ = moduleName
		_ = item
		if itemID != id {
			t.Fatalf("unexpected id: %s", itemID)
		}
		return want
	})

	m.Load(map[string]any{"id": id, "dsn": "x"})

	if Get(id) != want {
		t.Fatalf("engine mismatch: noerr")
	}
}
