package goetlstarter

import (
	"context"
	"database/sql"
	"errors"
	"time"

	gormsource "github.com/kordar/goetl-gorm/source"
	"github.com/kordar/goetl/checkpoint"
	"github.com/kordar/goetl/engine"
	"gorm.io/gorm"
)

type DBWalkerSpec interface {
	CheckpointKey() string
	PageSize() int
	MaxItems() int
	Store() checkpoint.CheckpointStore
	BuildQuery(ctx context.Context, cur checkpoint.Cursor, limit int) (string, []any, error)
	ExtractCursor(row map[string]any) (checkpoint.Cursor, error)
	MapRow(row map[string]any) (any, error)
	BuildAck(row map[string]any, cur checkpoint.Cursor) (string, error)
}

type DBWalkerSpecOptional interface {
	Name() string
	Partition() string
	BuildAttrs(row map[string]any, cur checkpoint.Cursor) (map[string]any, error)
}

// ErrNotImplemented 基础对象，便于后续通过嵌入扩展
var ErrNotImplemented = errors.New("not implemented")

type BaseDBWalkerSpec struct {
	CheckpointKeyValue string
	PageSizeValue      int
	MaxItemsValue      int
	NameValue          string
	PartitionValue     string
	StoreValue         checkpoint.CheckpointStore
}

func (b BaseDBWalkerSpec) CheckpointKey() string { return b.CheckpointKeyValue }
func (b BaseDBWalkerSpec) PageSize() int         { return b.PageSizeValue }
func (b BaseDBWalkerSpec) MaxItems() int         { return b.MaxItemsValue }
func (b BaseDBWalkerSpec) Name() string          { return b.NameValue }
func (b BaseDBWalkerSpec) Partition() string     { return b.PartitionValue }
func (b BaseDBWalkerSpec) Store() checkpoint.CheckpointStore {
	return b.StoreValue
}
func (b BaseDBWalkerSpec) BuildQuery(ctx context.Context, cur checkpoint.Cursor, limit int) (string, []any, error) {
	return "", nil, ErrNotImplemented
}
func (b BaseDBWalkerSpec) ExtractCursor(row map[string]any) (checkpoint.Cursor, error) {
	return checkpoint.Cursor{}, ErrNotImplemented
}
func (b BaseDBWalkerSpec) MapRow(row map[string]any) (any, error) { return row, nil }
func (b BaseDBWalkerSpec) BuildAck(row map[string]any, cur checkpoint.Cursor) (string, error) {
	return "", nil
}
func (b BaseDBWalkerSpec) BuildAttrs(row map[string]any, cur checkpoint.Cursor) (map[string]any, error) {
	return nil, nil
}

func LoadDBWalkerWithSQL(e *engine.Engine, db *sql.DB, spec DBWalkerSpec, interval, retry time.Duration, stopOnError bool) {
	w := &gormsource.DBWalker{
		SQL:           db,
		Store:         spec.Store(),
		CheckpointKey: spec.CheckpointKey(),
		PageSize:      spec.PageSize(),
		MaxItems:      spec.MaxItems(),
		BuildQuery:    spec.BuildQuery,
		ExtractCursor: spec.ExtractCursor,
		MapRow:        spec.MapRow,
		BuildAck:      spec.BuildAck,
	}
	if opt, ok := spec.(DBWalkerSpecOptional); ok {
		w.NameValue = opt.Name()
		w.Partition = opt.Partition()
		w.BuildAttrs = opt.BuildAttrs
	}
	tk := gormsource.NewDBWalkerTicker(w, interval, retry, stopOnError)
	e.LoadSource(tk)
}

func LoadDBWalkerWithGorm(e *engine.Engine, gdb *gorm.DB, spec DBWalkerSpec, interval, retry time.Duration, stopOnError bool) {
	w := &gormsource.DBWalker{
		Gorm:          gdb,
		Store:         spec.Store(),
		CheckpointKey: spec.CheckpointKey(),
		PageSize:      spec.PageSize(),
		MaxItems:      spec.MaxItems(),
		BuildQuery:    spec.BuildQuery,
		ExtractCursor: spec.ExtractCursor,
		MapRow:        spec.MapRow,
		BuildAck:      spec.BuildAck,
	}
	if opt, ok := spec.(DBWalkerSpecOptional); ok {
		w.NameValue = opt.Name()
		w.Partition = opt.Partition()
		w.BuildAttrs = opt.BuildAttrs
	}
	tk := gormsource.NewDBWalkerTicker(w, interval, retry, stopOnError)
	e.LoadSource(tk)
}
